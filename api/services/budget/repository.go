package budget

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
	stmt   *sqlx.Stmt
}

type budgetRow struct {
	Income      sql.NullInt64 `db:"income"`
	Expense     sql.NullInt64 `db:"expense"`
	LastUpdated sql.NullTime  `db:"last_updated"`
}

func initRepository(app *contracts.App) contracts.BudgetRepository {
	return &repository{
		reader: app.Ds.ReaderDB,
		stmt:   datasources.Prepare(app.Ds.ReaderDB, getMonthlyBudget),
	}
}

func (r *repository) GetMonthly(ctx context.Context, userID int64, year int, month int) (*entities.Budget, error) {
	row := budgetRow{}
	if err := r.stmt.GetContext(ctx, &row, userID, year, month); err != nil {
		return nil, err
	}

	b := &entities.Budget{
		Income:  row.Income.Int64,
		Expense: row.Expense.Int64,
	}
	if row.LastUpdated.Valid {
		b.LastUpdated = &row.LastUpdated.Time
	}
	return b, nil
}
