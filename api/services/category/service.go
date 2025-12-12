package category

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"finlog-api/api/contracts"
	"finlog-api/api/entities"
)

var (
	errCategoryNotFound = errors.New("category not found")
	errCategoryExists   = errors.New("category already exists")
	errInvalidCategory  = errors.New("invalid category input")
)

type Service struct {
	app  *contracts.App
	repo contracts.CategoryRepository
}

func Init(app *contracts.App) contracts.CategoryService {
	repo := initRepository(app)
	return &Service{app: app, repo: repo}
}

func (s *Service) ListCategories(ctx context.Context, userID int64, filter contracts.CategoryFilter) ([]entities.Category, error) {
	return s.repo.List(ctx, userID, filter)
}

func (s *Service) CreateCategory(ctx context.Context, userID int64, name string, isExpense bool, iconKey string) (*entities.Category, error) {
	if err := validateCategoryInput(name); err != nil {
		return nil, err
	}
	if iconKey == "" {
		iconKey = "category"
	}

	existing, err := s.repo.FindByName(ctx, name, isExpense, userID)
	if err == nil {
		if existing.IsActive {
			return nil, errCategoryExists
		}
		existing.Name = name
		existing.IconKey = iconKey
		existing.IsExpense = isExpense
		existing.IsActive = true
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	category := &entities.Category{
		UserID:    userID,
		Name:      name,
		IsExpense: isExpense,
		IconKey:   iconKey,
	}
	id, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	category.ID = id
	return category, nil
}

func (s *Service) UpdateCategory(ctx context.Context, userID, categoryID int64, name string, isExpense bool, iconKey string) error {
	if categoryID <= 0 {
		return errInvalidCategory
	}
	if err := validateCategoryInput(name); err != nil {
		return err
	}
	if iconKey == "" {
		iconKey = "category"
	}

	if existing, err := s.repo.FindByName(ctx, name, isExpense, userID); err == nil && existing.ID != categoryID && existing.IsActive {
		return errCategoryExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	category := &entities.Category{
		ID:        categoryID,
		UserID:    userID,
		Name:      name,
		IsExpense: isExpense,
		IconKey:   iconKey,
		IsActive:  true,
	}
	if err := s.repo.Update(ctx, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errCategoryNotFound
		}
		return err
	}
	return nil
}

func (s *Service) DeleteCategory(ctx context.Context, userID, categoryID int64) error {
	if categoryID <= 0 {
		return errInvalidCategory
	}
	if err := s.repo.Delete(ctx, categoryID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errCategoryNotFound
		}
		return err
	}
	return nil
}

func validateCategoryInput(name string) error {
	if strings.TrimSpace(name) == "" {
		return errInvalidCategory
	}
	if len(name) > 100 {
		return errors.New("category name too long")
	}
	return nil
}

// Exported errors for handlers.
func ErrCategoryNotFound() error { return errCategoryNotFound }
func ErrCategoryExists() error   { return errCategoryExists }
func ErrInvalidCategory() error  { return errInvalidCategory }
