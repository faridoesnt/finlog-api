package contracts

import (
	"context"

	"finlog-api/api/entities"
	"finlog-api/api/models/request"
)

type ImportRepository interface {
	InsertBatch(ctx context.Context, userID int64, items []entities.ImportedTransaction) error
	ListBatches(ctx context.Context, userID int64) ([]entities.ImportBatch, error)
	DeleteBatch(ctx context.Context, userID, batchID int64) (int64, error)
}

type ImportService interface {
	StoreBatch(ctx context.Context, userID int64, payload request.ImportBatchRequest) error
	ListHistory(ctx context.Context, userID int64) ([]entities.ImportBatch, error)
	UndoBatch(ctx context.Context, userID, batchID int64) (int64, error)
}
