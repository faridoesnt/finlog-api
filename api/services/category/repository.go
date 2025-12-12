package category

import (
	"context"
	"database/sql"

	"finlog-api/api/contracts"
	"finlog-api/api/datasources"
	"finlog-api/api/entities"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	reader *sqlx.DB
	writer *sqlx.DB
	stmt   struct {
		findByID   *sqlx.Stmt
		findByName *sqlx.Stmt
		insert     *sqlx.Stmt
		update     *sqlx.Stmt
		delete     *sqlx.Stmt
	}
}

func initRepository(app *contracts.App) contracts.CategoryRepository {
	return &repository{
		reader: app.Ds.ReaderDB,
		writer: app.Ds.WriterDB,
		stmt: struct {
			findByID   *sqlx.Stmt
			findByName *sqlx.Stmt
			insert     *sqlx.Stmt
			update     *sqlx.Stmt
			delete     *sqlx.Stmt
		}{
			findByID:   datasources.Prepare(app.Ds.ReaderDB, findCategoryByID),
			findByName: datasources.Prepare(app.Ds.ReaderDB, findCategoryByName),
			insert:     datasources.Prepare(app.Ds.WriterDB, insertCategory),
			update:     datasources.Prepare(app.Ds.WriterDB, updateCategory),
			delete:     datasources.Prepare(app.Ds.WriterDB, deleteCategory),
		},
	}
}

func NewRepository(app *contracts.App) contracts.CategoryRepository {
	return initRepository(app)
}
func (r *repository) List(ctx context.Context, userID int64, filter contracts.CategoryFilter) ([]entities.Category, error) {
	query := listCategories
	args := []interface{}{userID}
	if filter.IsExpense != nil {
		query += " AND is_expense = ?"
		args = append(args, *filter.IsExpense)
	}
	query += " ORDER BY name ASC"

	var categories []entities.Category
	if err := r.reader.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) FindByID(ctx context.Context, id, userID int64) (*entities.Category, error) {
	cat := new(entities.Category)
	if err := r.stmt.findByID.GetContext(ctx, cat, id, userID); err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *repository) FindByName(ctx context.Context, name string, isExpense bool, userID int64) (*entities.Category, error) {
	cat := new(entities.Category)
	if err := r.stmt.findByName.GetContext(ctx, cat, userID, name, isExpense); err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *repository) Create(ctx context.Context, category *entities.Category) (int64, error) {
	res, err := r.stmt.insert.ExecContext(ctx, category.UserID, category.Name, category.IsExpense, category.IconKey)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *repository) Update(ctx context.Context, category *entities.Category) error {
	res, err := r.stmt.update.ExecContext(ctx, category.Name, category.IsExpense, category.IconKey, boolToInt(category.IsActive), category.ID, category.UserID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id, userID int64) error {
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

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
