package entities

import "time"

type ImportBatch struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	BatchSize int       `db:"batch_size"`
	CreatedAt time.Time `db:"created_at"`
}

type ImportedTransaction struct {
	Ciphertext string
	Nonce      string
	Tag        string
	OccurredAt time.Time
	IsExpense  bool
	CategoryID int64
}
