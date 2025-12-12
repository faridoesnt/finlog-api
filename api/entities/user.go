package entities

import "time"

// User represents user.
type User struct {
	ID                    int64      `db:"id" json:"id"`
	Email                 string     `db:"email" json:"email"`
	Name                  string     `db:"name" json:"name"`
	Role                  string     `db:"role" json:"role"`
	Password              string     `db:"password" json:"-"`
	IsVerified            bool       `db:"is_verified" json:"is_verified"`
	VerificationToken     *string    `db:"verification_token" json:"-"`
	VerificationExpiresAt *time.Time `db:"verification_expires_at" json:"-"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}
