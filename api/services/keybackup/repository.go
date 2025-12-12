package keybackup

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"
)

type repository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func initRepository(app *contracts.App) contracts.KeyBackupRepository {
	return &repository{
		reader: app.Ds.ReaderDB,
		writer: app.Ds.WriterDB,
	}
}

func (r *repository) GetActive(ctx context.Context, userID int64) (*entities.UserEncryptedDataKey, error) {
	key := new(entities.UserEncryptedDataKey)
	if err := r.reader.GetContext(ctx, key, selectActiveKeyQuery, userID); err != nil {
		return nil, err
	}
	return key, nil
}

func (r *repository) Insert(ctx context.Context, exec sqlx.ExtContext, key *entities.UserEncryptedDataKey) (int64, error) {
	res, err := exec.ExecContext(ctx, insertKeyQuery, key.UserID, key.EncryptedDataKey, key.Salt, boolToInt(key.IsActive))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *repository) DeactivateActive(ctx context.Context, exec sqlx.ExtContext, userID int64, rotatedAt time.Time) (int64, error) {
	res, err := exec.ExecContext(ctx, deactivateActiveKeyQuery, rotatedAt, rotatedAt, userID)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, nil
}

func (r *repository) RotationSummary(ctx context.Context, userID int64) (int64, *time.Time, error) {
	summary := struct {
		RotationCount int64      `db:"rotation_count"`
		LastRotatedAt *time.Time `db:"last_rotated_at"`
	}{}
	if err := r.reader.GetContext(ctx, &summary, rotationSummaryQuery, userID); err != nil {
		return 0, nil, err
	}
	return summary.RotationCount, summary.LastRotatedAt, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
