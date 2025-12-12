package importbatch

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/entities"
	"finlog-api/api/models/request"
	"finlog-api/api/services/category"
)

var (
	ErrInvalidImportInput    = errors.New("invalid import payload")
	ErrRateLimitExceeded     = errors.New("import rate limit exceeded")
	ErrUndoRateLimitExceeded = errors.New("undo rate limit exceeded")
	ErrImportBatchNotFound   = errors.New("import batch not found")
)

type Service struct {
	app          *contracts.App
	repo         contracts.ImportRepository
	limiter      *rateLimiter
	undoLimiter  *rateLimiter
	categoryRepo contracts.CategoryRepository
}

func Init(app *contracts.App) contracts.ImportService {
	limit, window := parseRateLimit(app.Config)
	undoLimit, undoWindow := parseUndoRateLimit(app.Config)
	return &Service{
		app:          app,
		repo:         initRepository(app),
		limiter:      newRateLimiter(limit, window),
		undoLimiter:  newRateLimiter(undoLimit, undoWindow),
		categoryRepo: category.NewRepository(app),
	}
}

func (s *Service) StoreBatch(ctx context.Context, userID int64, payload request.ImportBatchRequest) error {
	if len(payload.Items) == 0 {
		return ErrInvalidImportInput
	}
	if !s.limiter.Allow(userID) {
		return ErrRateLimitExceeded
	}

	items := make([]entities.ImportedTransaction, 0, len(payload.Items))
	categoryCache := make(map[int64]*entities.Category)
	for _, item := range payload.Items {
		if strings.TrimSpace(item.Ciphertext) == "" ||
			strings.TrimSpace(item.Nonce) == "" ||
			strings.TrimSpace(item.Tag) == "" {
			return ErrInvalidImportInput
		}
		occurredAtRaw := strings.TrimSpace(item.OccurredAt)
		if occurredAtRaw == "" {
			return ErrInvalidImportInput
		}
		occurredAt, err := time.Parse(time.RFC3339, occurredAtRaw)
		if err != nil {
			occurredAt, err = time.Parse(time.RFC3339Nano, occurredAtRaw)
		}
		if err != nil {
			return ErrInvalidImportInput
		}
		if item.CategoryID <= 0 {
			return ErrInvalidImportInput
		}

		category, ok := categoryCache[item.CategoryID]
		if !ok {
			category, err = s.categoryRepo.FindByID(ctx, item.CategoryID, userID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return ErrInvalidImportInput
				}
				return err
			}
			categoryCache[item.CategoryID] = category
		}
		if category == nil {
			return ErrInvalidImportInput
		}
		if category.IsExpense != item.IsExpense {
			return ErrInvalidImportInput
		}

		items = append(items, entities.ImportedTransaction{
			Ciphertext: item.Ciphertext,
			Nonce:      item.Nonce,
			Tag:        item.Tag,
			OccurredAt: occurredAt,
			IsExpense:  item.IsExpense,
			CategoryID: item.CategoryID,
		})
	}

	if err := s.repo.InsertBatch(ctx, userID, items); err != nil {
		return err
	}

	s.app.Logger.Info().
		Int64("user_id", userID).
		Int("batch_size", len(items)).
		Msg("Import batch stored")

	return nil
}

func (s *Service) ListHistory(ctx context.Context, userID int64) ([]entities.ImportBatch, error) {
	return s.repo.ListBatches(ctx, userID)
}

func (s *Service) UndoBatch(ctx context.Context, userID, batchID int64) (int64, error) {
	if batchID <= 0 {
		return 0, ErrInvalidImportInput
	}
	if !s.undoLimiter.Allow(userID) {
		return 0, ErrUndoRateLimitExceeded
	}

	deleted, err := s.repo.DeleteBatch(ctx, userID, batchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrImportBatchNotFound
		}
		return 0, err
	}

	s.app.Logger.Info().
		Int64("user_id", userID).
		Int64("batch_id", batchID).
		Int64("deleted_count", deleted).
		Msg("Import batch undone")

	return deleted, nil
}

func parseRateLimit(config map[string]string) (int, time.Duration) {
	limit := 5
	if raw := config[constants.ImportRateLimitBatches]; raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	window := time.Minute
	if raw := config[constants.ImportRateLimitWindow]; raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			window = parsed
		}
	}
	return limit, window
}

func parseUndoRateLimit(config map[string]string) (int, time.Duration) {
	limit := 3
	if raw := config[constants.ImportUndoRateLimitRequests]; raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	window := time.Minute
	if raw := config[constants.ImportUndoRateLimitWindow]; raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			window = parsed
		}
	}
	return limit, window
}

type rateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	counters map[int64]*rateCounter
}

type rateCounter struct {
	count       int
	windowStart time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &rateLimiter{
		window:   window,
		max:      limit,
		counters: make(map[int64]*rateCounter),
	}
}

func (l *rateLimiter) Allow(userID int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	counter, exists := l.counters[userID]
	if !exists || now.Sub(counter.windowStart) >= l.window {
		l.counters[userID] = &rateCounter{count: 1, windowStart: now}
		return true
	}
	if counter.count >= l.max {
		return false
	}
	counter.count++
	return true
}
