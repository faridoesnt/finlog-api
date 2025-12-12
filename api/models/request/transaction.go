package request

import "time"

type CreateTransaction struct {
	Ciphertext string    `json:"ciphertext" validate:"required"`
	Nonce      string    `json:"nonce" validate:"required"`
	Tag        string    `json:"tag" validate:"required"`
	Date       string    `json:"date" validate:"required"`
	OccurredAt time.Time `json:"-" validate:"-"`
	IsExpense  bool      `json:"isExpense" validate:"required"`
	Category   string    `json:"category" validate:"required"`
}
