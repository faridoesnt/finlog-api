package transaction

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"
	"finlog-api/api/models/request"
	"finlog-api/api/services/category"
)

const defaultRecentLimit = 10

var (
	errTransactionNotFound = errors.New("transaction not found")
	errInvalidTransaction  = errors.New("invalid transaction input")
	errCategoryNotFound    = errors.New("category not found")
	errUnsupportedBulk     = errors.New("bulk update not supported for encrypted payloads")
)

type Service struct {
	app          *contracts.App
	txRepo       contracts.TransactionRepository
	categoryRepo contracts.CategoryRepository
}

func Init(app *contracts.App) contracts.TransactionService {
	return &Service{
		app:          app,
		txRepo:       initRepository(app),
		categoryRepo: category.NewRepository(app),
	}
}

func (s *Service) GetTransactions(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error) {
	return s.txRepo.List(ctx, userID, year, month)
}

func (s *Service) GetRecentTransactions(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error) {
	return s.txRepo.ListRecent(ctx, userID, year, month, defaultRecentLimit)
}

func (s *Service) CreateTransaction(ctx context.Context, userID int64, input request.CreateTransaction) (*entities.Transaction, error) {
	if err := validateTransactionInput(input); err != nil {
		return nil, err
	}

	cat, err := s.resolveCategory(ctx, userID, input.Category, input.IsExpense)
	if err != nil {
		return nil, err
	}

	tx := &entities.Transaction{
		UserID:     userID,
		CategoryID: cat.ID,
		Category:   cat.Name,
		Ciphertext: input.Ciphertext,
		Nonce:      input.Nonce,
		Tag:        input.Tag,
		OccurredAt: input.OccurredAt,
		IsExpense:  input.IsExpense,
	}
	id, err := s.txRepo.Create(ctx, tx)
	if err != nil {
		return nil, err
	}
	tx.ID = id
	return tx, nil
}

func (s *Service) UpdateTransaction(ctx context.Context, userID int64, id int64, input request.CreateTransaction) error {
	if id <= 0 {
		return errInvalidTransaction
	}
	if err := validateTransactionInput(input); err != nil {
		return err
	}
	cat, err := s.resolveCategory(ctx, userID, input.Category, input.IsExpense)
	if err != nil {
		return err
	}

	tx := &entities.Transaction{
		ID:         id,
		UserID:     userID,
		CategoryID: cat.ID,
		Ciphertext: input.Ciphertext,
		Nonce:      input.Nonce,
		Tag:        input.Tag,
		OccurredAt: input.OccurredAt,
		IsExpense:  input.IsExpense,
	}
	if err := s.txRepo.Update(ctx, tx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *Service) UpdateNotes(ctx context.Context, userID int64, ids []int64, notes string) error {
	return errUnsupportedBulk
}

func (s *Service) UpdateAmount(ctx context.Context, userID int64, ids []int64, amount int64) error {
	return errUnsupportedBulk
}

func (s *Service) UpdateDate(ctx context.Context, userID int64, ids []int64, date time.Time) error {
	return errUnsupportedBulk
}

func (s *Service) DeleteTransaction(ctx context.Context, userID int64, id int64) error {
	if id <= 0 {
		return errInvalidTransaction
	}
	if err := s.txRepo.Delete(ctx, userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *Service) DeleteTransactions(ctx context.Context, userID int64, ids []int64) error {
	if len(ids) == 0 {
		return errInvalidTransaction
	}
	if err := s.txRepo.BulkDelete(ctx, userID, ids); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errTransactionNotFound
		}
		return err
	}
	return nil
}

func (s *Service) resolveCategory(ctx context.Context, userID int64, identifier string, isExpense bool) (*entities.Category, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, errCategoryNotFound
	}
	if id, err := strconv.ParseInt(identifier, 10, 64); err == nil && id > 0 {
		cat, err := s.categoryRepo.FindByID(ctx, id, userID)
		if err == nil && cat != nil {
			if cat.IsExpense != isExpense {
				return nil, errInvalidTransaction
			}
			return cat, nil
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	cat, err := s.categoryRepo.FindByName(ctx, identifier, isExpense, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errCategoryNotFound
		}
		return nil, err
	}
	if cat.IsExpense != isExpense {
		return nil, errInvalidTransaction
	}
	return cat, nil
}

func validateTransactionInput(input request.CreateTransaction) error {
	if input.OccurredAt.IsZero() {
		return errInvalidTransaction
	}
	if strings.TrimSpace(input.Ciphertext) == "" || strings.TrimSpace(input.Nonce) == "" || strings.TrimSpace(input.Tag) == "" {
		return errInvalidTransaction
	}
	return nil
}

// Exported errors for handlers.
func ErrTransactionNotFound() error { return errTransactionNotFound }
func ErrInvalidTransaction() error  { return errInvalidTransaction }
func ErrCategoryNotFound() error    { return errCategoryNotFound }
func ErrUnsupportedBulk() error     { return errUnsupportedBulk }
