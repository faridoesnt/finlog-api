package contracts

import (
	"context"

	"finlog-api/api/entities"
)

type BudgetRepository interface {
	GetMonthly(ctx context.Context, userID int64, year int, month int) (*entities.Budget, error)
}

type BudgetService interface {
	GetMonthly(ctx context.Context, userID int64, year int, month int) (*entities.Budget, error)
}
