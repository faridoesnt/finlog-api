package importbatch

import (
	"context"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"

	"github.com/jmoiron/sqlx"
)

const (
	insertImportBatchSQL = `
		INSERT INTO import_batches (user_id, batch_size)
		VALUES (?, ?)
	`
	insertTransactionSQL = `
		INSERT INTO transactions (user_id, category_id, payload_ciphertext, payload_nonce, payload_tag, occurred_at, is_expense, batch_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	listImportBatchesSQL = `
		SELECT id, user_id, batch_size, created_at
		FROM import_batches
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	selectImportBatchSQL = `
		SELECT id
		FROM import_batches
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`
	deleteTransactionsByBatchSQL = `
		DELETE FROM transactions
		WHERE user_id = ? AND batch_id = ?
	`
	deleteImportBatchSQL = `
		DELETE FROM import_batches
		WHERE id = ? AND user_id = ?
	`
)

type repository struct {
	writer *sqlx.DB
}

func initRepository(app *contracts.App) contracts.ImportRepository {
	return &repository{
		writer: app.Ds.WriterDB,
	}
}

func (r *repository) InsertBatch(
	ctx context.Context,
	userID int64,
	items []entities.ImportedTransaction,
) error {
	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, insertImportBatchSQL, userID, len(items))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	batchID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, item := range items {
		if _, err := tx.ExecContext(
			ctx,
			insertTransactionSQL,
			userID,
			item.CategoryID,
			item.Ciphertext,
			item.Nonce,
			item.Tag,
			item.OccurredAt,
			item.IsExpense,
			batchID,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *repository) ListBatches(
	ctx context.Context,
	userID int64,
) ([]entities.ImportBatch, error) {
	var batches []entities.ImportBatch
	if err := r.writer.SelectContext(ctx, &batches, listImportBatchesSQL, userID); err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *repository) DeleteBatch(
	ctx context.Context,
	userID,
	batchID int64,
) (int64, error) {
	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var exists int
	if err := tx.GetContext(ctx, &exists, selectImportBatchSQL, batchID, userID); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	result, err := tx.ExecContext(ctx, deleteTransactionsByBatchSQL, userID, batchID)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, deleteImportBatchSQL, batchID, userID); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return deleted, nil
}
