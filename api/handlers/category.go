package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"finlog-api/api/contracts"
	"finlog-api/api/models/responses"
	"finlog-api/api/services/category"
)

// GetCategories returns categories for the current user.
func GetCategories(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	filter := contracts.CategoryFilter{}
	switch strings.ToLower(c.Query("type")) {
	case "expense":
		filter.IsExpense = boolPtr(true)
	case "income":
		filter.IsExpense = boolPtr(false)
	}

	categories, err := app.Services.Categories.ListCategories(context.Background(), userID, filter)
	if err != nil {
		return responses.InternalServerError(err)
	}
	return c.JSON(categories)
}

// CreateCategory adds a new category for the user.
func CreateCategory(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	type req struct {
		Name      string `json:"name"`
		IsExpense bool   `json:"isExpense"`
		IconKey   string `json:"icon"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}

	cat, err := app.Services.Categories.CreateCategory(context.Background(), userID, body.Name, body.IsExpense, body.IconKey)
	if err != nil {
		switch {
		case errors.Is(err, category.ErrCategoryExists()):
			return responses.Conflict(err)
		case errors.Is(err, category.ErrInvalidCategory()):
			return responses.BadRequest(err)
		default:
			return responses.BadRequest(err)
		}
	}
	return c.Status(fiber.StatusCreated).JSON(cat)
}

// UpdateCategory updates category attributes.
func UpdateCategory(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return responses.BadRequest(errors.New("invalid category id"))
	}
	type req struct {
		Name      string `json:"name"`
		IsExpense bool   `json:"isExpense"`
		IconKey   string `json:"icon"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	if err := app.Services.Categories.UpdateCategory(context.Background(), userID, categoryID, body.Name, body.IsExpense, body.IconKey); err != nil {
		switch {
		case errors.Is(err, category.ErrCategoryExists()):
			return responses.Conflict(err)
		case errors.Is(err, category.ErrCategoryNotFound()):
			return responses.NotFound(err)
		default:
			return responses.BadRequest(err)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteCategory removes a category.
func DeleteCategory(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return responses.BadRequest(errors.New("invalid category id"))
	}
	if err := app.Services.Categories.DeleteCategory(context.Background(), userID, categoryID); err != nil {
		if errors.Is(err, category.ErrCategoryNotFound()) {
			return responses.NotFound(err)
		}
		return responses.BadRequest(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func boolPtr(v bool) *bool {
	return &v
}
