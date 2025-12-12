package contracts

import (
	"context"
	"time"

	"finlog-api/api/entities"
	"finlog-api/api/models/request"
)

type TransactionRepository interface {
	List(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error)
	ListRecent(ctx context.Context, userID int64, year int, month int, limit int) ([]entities.Transaction, error)
	Create(ctx context.Context, tx *entities.Transaction) (int64, error)
	Update(ctx context.Context, tx *entities.Transaction) error
	BulkUpdateNotes(ctx context.Context, userID int64, ids []int64, notes string) error
	BulkUpdateAmount(ctx context.Context, userID int64, ids []int64, amount int64) error
	BulkUpdateDate(ctx context.Context, userID int64, ids []int64, date time.Time) error
	Delete(ctx context.Context, userID int64, id int64) error
	BulkDelete(ctx context.Context, userID int64, ids []int64) error
	FindByID(ctx context.Context, id, userID int64) (*entities.Transaction, error)
}

type TransactionService interface {
	GetTransactions(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error)
	GetRecentTransactions(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error)
	CreateTransaction(ctx context.Context, userID int64, input request.CreateTransaction) (*entities.Transaction, error)
	UpdateTransaction(ctx context.Context, userID int64, id int64, input request.CreateTransaction) error
	UpdateNotes(ctx context.Context, userID int64, ids []int64, notes string) error
	UpdateAmount(ctx context.Context, userID int64, ids []int64, amount int64) error
	UpdateDate(ctx context.Context, userID int64, ids []int64, date time.Time) error
	DeleteTransaction(ctx context.Context, userID int64, id int64) error
	DeleteTransactions(ctx context.Context, userID int64, ids []int64) error
}
