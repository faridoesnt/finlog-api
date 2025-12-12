package contracts

import (
	"context"
	"time"

	"finlog-api/api/entities"

	"github.com/jmoiron/sqlx"
)

// KeyBackupStatus describes what metadata the backend can return to the client.
type KeyBackupStatus struct {
	HasActiveKey  bool       `json:"has_active_key"`
	RotationCount int64      `json:"rotation_count"`
	LastRotatedAt *time.Time `json:"last_rotated_at,omitempty"`
}

// KeyBackupRepository provides low-level access to encrypted key rows.
type KeyBackupRepository interface {
	GetActive(ctx context.Context, userID int64) (*entities.UserEncryptedDataKey, error)
	Insert(ctx context.Context, exec sqlx.ExtContext, key *entities.UserEncryptedDataKey) (int64, error)
	DeactivateActive(ctx context.Context, exec sqlx.ExtContext, userID int64, rotatedAt time.Time) (int64, error)
	RotationSummary(ctx context.Context, userID int64) (int64, *time.Time, error)
}

// KeyBackupService coordinates business logic around encrypted keys.
type KeyBackupService interface {
	StoreKeyBackup(ctx context.Context, userID int64, encryptedKey, salt string) (*entities.UserEncryptedDataKey, error)
	RotateKey(ctx context.Context, userID int64, encryptedKey, salt string) (*entities.UserEncryptedDataKey, error)
	GetActiveKey(ctx context.Context, userID int64) (*entities.UserEncryptedDataKey, error)
	GetKeyStatus(ctx context.Context, userID int64) (*KeyBackupStatus, error)
}
