package entities

import "time"

// Transaction represents a single income/expense record.
type Transaction struct {
	ID         int64     `db:"id" json:"id"`
	UserID     int64     `db:"user_id" json:"-"`
	CategoryID int64     `db:"category_id" json:"category_id,omitempty"`
	Category   string    `db:"category_name" json:"category"`
	Ciphertext string    `db:"payload_ciphertext" json:"ciphertext"`
	Nonce      string    `db:"payload_nonce" json:"nonce"`
	Tag        string    `db:"payload_tag" json:"tag"`
	OccurredAt time.Time `db:"occurred_at" json:"date"`
	IsExpense  bool      `db:"is_expense" json:"isExpense"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
