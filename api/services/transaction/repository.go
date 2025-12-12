package transaction

import (
	"context"
	"database/sql"
	"time"

	"finlog-api/api/contracts"
	"finlog-api/api/datasources"
	"finlog-api/api/entities"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	reader *sqlx.DB
	writer *sqlx.DB
	stmt   struct {
		findByID *sqlx.Stmt
		insert   *sqlx.Stmt
		update   *sqlx.Stmt
		delete   *sqlx.Stmt
	}
}

func initRepository(app *contracts.App) contracts.TransactionRepository {
	return &repository{
		reader: app.Ds.ReaderDB,
		writer: app.Ds.WriterDB,
		stmt: struct {
			findByID *sqlx.Stmt
			insert   *sqlx.Stmt
			update   *sqlx.Stmt
			delete   *sqlx.Stmt
		}{
			findByID: datasources.Prepare(app.Ds.ReaderDB, findTransactionByID),
			insert:   datasources.Prepare(app.Ds.WriterDB, insertTransaction),
			update:   datasources.Prepare(app.Ds.WriterDB, updateTransaction),
			delete:   datasources.Prepare(app.Ds.WriterDB, deleteTransaction),
		},
	}
}

func (r *repository) List(ctx context.Context, userID int64, year int, month int) ([]entities.Transaction, error) {
	var txs []entities.Transaction
	if err := r.reader.SelectContext(ctx, &txs, listTransactions, userID, year, month); err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *repository) ListRecent(ctx context.Context, userID int64, year int, month int, limit int) ([]entities.Transaction, error) {
	query := listTransactions + " LIMIT ?"
	args := []interface{}{userID, year, month, limit}

	var txs []entities.Transaction
	if err := r.reader.SelectContext(ctx, &txs, query, args...); err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *repository) FindByID(ctx context.Context, id, userID int64) (*entities.Transaction, error) {
	tx := new(entities.Transaction)
	if err := r.stmt.findByID.GetContext(ctx, tx, id, userID); err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *repository) Create(ctx context.Context, tx *entities.Transaction) (int64, error) {
	res, err := r.stmt.insert.ExecContext(
		ctx,
		tx.UserID,
		tx.CategoryID,
		tx.Ciphertext,
		tx.Nonce,
		tx.Tag,
		tx.OccurredAt,
		tx.IsExpense,
		nil,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *repository) Update(ctx context.Context, tx *entities.Transaction) error {
	res, err := r.stmt.update.ExecContext(ctx, tx.Ciphertext, tx.Nonce, tx.Tag, tx.OccurredAt, tx.IsExpense, tx.CategoryID, tx.ID, tx.UserID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) BulkUpdateNotes(ctx context.Context, userID int64, ids []int64, notes string) error {
	return r.execBulk(ctx, userID, ids, "payload_ciphertext = payload_ciphertext", nil)
}

func (r *repository) BulkUpdateAmount(ctx context.Context, userID int64, ids []int64, amount int64) error {
	return r.execBulk(ctx, userID, ids, "payload_ciphertext = payload_ciphertext", nil)
}

func (r *repository) BulkUpdateDate(ctx context.Context, userID int64, ids []int64, date time.Time) error {
	return r.execBulk(ctx, userID, ids, "occurred_at = ?", date)
}

func (r *repository) Delete(ctx context.Context, userID int64, id int64) error {
	res, err := r.stmt.delete.ExecContext(ctx, id, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) BulkDelete(ctx context.Context, userID int64, ids []int64) error {
	if len(ids) == 0 {
		return sql.ErrNoRows
	}
	query, args, err := sqlx.In("DELETE FROM transactions WHERE user_id = ? AND id IN (?)", userID, ids)
	if err != nil {
		return err
	}
	query = r.writer.Rebind(query)
	res, err := r.writer.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) execBulk(ctx context.Context, userID int64, ids []int64, setClause string, value interface{}) error {
	if len(ids) == 0 {
		return sql.ErrNoRows
	}
	query := "UPDATE transactions SET " + setClause + ", updated_at = NOW() WHERE user_id = ? AND id IN (?)"
	params := []interface{}{value, userID, ids}

	q, args, err := sqlx.In(query, params...)
	if err != nil {
		return err
	}
	q = r.writer.Rebind(q)
	res, err := r.writer.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
