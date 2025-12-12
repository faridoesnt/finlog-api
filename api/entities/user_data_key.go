package entities

import "time"

// UserEncryptedDataKey stores encrypted data key metadata per user while keeping actual key material opaque.
type UserEncryptedDataKey struct {
	ID               int64      `db:"id" json:"id"`
	UserID           int64      `db:"user_id" json:"-"`
	EncryptedDataKey string     `db:"encrypted_data_key" json:"encrypted_data_key"`
	Salt             string     `db:"salt" json:"salt"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	RotatedAt        *time.Time `db:"rotated_at" json:"rotated_at,omitempty"`
	DeletedAt        *time.Time `db:"deleted_at" json:"-"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}
