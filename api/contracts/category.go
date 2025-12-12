package contracts

import (
	"context"

	"finlog-api/api/entities"
)

type CategoryFilter struct {
	IsExpense *bool
}

type CategoryRepository interface {
	List(ctx context.Context, userID int64, filter CategoryFilter) ([]entities.Category, error)
	Create(ctx context.Context, category *entities.Category) (int64, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id, userID int64) error
	FindByID(ctx context.Context, id, userID int64) (*entities.Category, error)
	FindByName(ctx context.Context, name string, isExpense bool, userID int64) (*entities.Category, error)
}

type CategoryService interface {
	ListCategories(ctx context.Context, userID int64, filter CategoryFilter) ([]entities.Category, error)
	CreateCategory(ctx context.Context, userID int64, name string, isExpense bool, iconKey string) (*entities.Category, error)
	UpdateCategory(ctx context.Context, userID, categoryID int64, name string, isExpense bool, iconKey string) error
	DeleteCategory(ctx context.Context, userID, categoryID int64) error
}
