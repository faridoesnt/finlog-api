package keybackup

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"
)

const maxEncryptedPayload = 8192

var (
	errInvalidPayload  = errors.New("encrypted key payload is invalid")
	errActiveKeyExists = errors.New("an active encrypted key already exists")
	errKeyNotFound     = errors.New("active encrypted key not found")
)

type Service struct {
	app  *contracts.App
	repo contracts.KeyBackupRepository
}

func Init(app *contracts.App) contracts.KeyBackupService {
	return &Service{
		app:  app,
		repo: initRepository(app),
	}
}

func (s *Service) StoreKeyBackup(ctx context.Context, userID int64, encryptedKey, salt string) (*entities.UserEncryptedDataKey, error) {
	if err := validatePayload(encryptedKey, salt); err != nil {
		return nil, err
	}

	if key, err := s.repo.GetActive(ctx, userID); err == nil && key != nil {
		return nil, errActiveKeyExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	tx, err := s.app.Ds.WriterDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	newKey := &entities.UserEncryptedDataKey{
		UserID:           userID,
		EncryptedDataKey: encryptedKey,
		Salt:             salt,
		IsActive:         true,
	}

	id, err := s.repo.Insert(ctx, tx, newKey)
	if err != nil {
		return nil, err
	}
	newKey.ID = id

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	tx = nil

	s.logAction(userID, "backup_created")
	return newKey, nil
}

func (s *Service) RotateKey(ctx context.Context, userID int64, encryptedKey, salt string) (*entities.UserEncryptedDataKey, error) {
	if err := validatePayload(encryptedKey, salt); err != nil {
		return nil, err
	}

	current, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errKeyNotFound
		}
		return nil, err
	}
	if current == nil {
		return nil, errKeyNotFound
	}

	tx, err := s.app.Ds.WriterDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	now := time.Now().UTC()
	affected, err := s.repo.DeactivateActive(ctx, tx, userID, now)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errKeyNotFound
	}

	newKey := &entities.UserEncryptedDataKey{
		UserID:           userID,
		EncryptedDataKey: encryptedKey,
		Salt:             salt,
		IsActive:         true,
	}
	id, err := s.repo.Insert(ctx, tx, newKey)
	if err != nil {
		return nil, err
	}
	newKey.ID = id

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	tx = nil

	s.logAction(userID, "backup_rotated")
	return newKey, nil
}

func (s *Service) GetActiveKey(ctx context.Context, userID int64) (*entities.UserEncryptedDataKey, error) {
	return s.repo.GetActive(ctx, userID)
}

func (s *Service) GetKeyStatus(ctx context.Context, userID int64) (*contracts.KeyBackupStatus, error) {
	hasActive := true
	if _, err := s.repo.GetActive(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hasActive = false
		} else {
			return nil, err
		}
	}

	count, lastRotatedAt, err := s.repo.RotationSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &contracts.KeyBackupStatus{
		HasActiveKey:  hasActive,
		RotationCount: count,
		LastRotatedAt: lastRotatedAt,
	}, nil
}

func (s *Service) logAction(userID int64, action string) {
	if s.app == nil || s.app.Logger == nil {
		return
	}
	s.app.Logger.Info().
		Int64("user_id", userID).
		Str("action", action).
		Msg("key backup event")
}

func validatePayload(encryptedKey, salt string) error {
	if strings.TrimSpace(encryptedKey) == "" || strings.TrimSpace(salt) == "" {
		return errInvalidPayload
	}
	if len(encryptedKey) > maxEncryptedPayload || len(salt) > maxEncryptedPayload {
		return errInvalidPayload
	}
	return nil
}

func ErrInvalidPayload() error  { return errInvalidPayload }
func ErrActiveKeyExists() error { return errActiveKeyExists }
func ErrKeyNotFound() error     { return errKeyNotFound }
