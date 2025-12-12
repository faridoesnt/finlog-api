package entities

import "time"

// Category represents a spend/income classification owned by a user.
type Category struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"-"`
	Name      string    `db:"name" json:"name"`
	IsExpense bool      `db:"is_expense" json:"isExpense"`
	IconKey   string    `db:"icon_key" json:"icon"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
