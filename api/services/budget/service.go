package budget

import (
	"context"
	"errors"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"
)

var errInvalidPeriod = errors.New("invalid period")

type Service struct {
	app  *contracts.App
	repo contracts.BudgetRepository
}

func Init(app *contracts.App) contracts.BudgetService {
	return &Service{
		app:  app,
		repo: initRepository(app),
	}
}

func (s *Service) GetMonthly(ctx context.Context, userID int64, year int, month int) (*entities.Budget, error) {
	if year <= 0 || month < 1 || month > 12 {
		return nil, errInvalidPeriod
	}
	return s.repo.GetMonthly(ctx, userID, year, month)
}

func ErrInvalidPeriod() error { return errInvalidPeriod }
