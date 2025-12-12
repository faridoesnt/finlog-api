package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"

	"finlog-api/api/models/responses"
	"finlog-api/api/services/budget"
)

// GetBudget returns aggregated income/expense for the period.
func GetBudget(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	year, month, err := parsePeriod(c)
	if err != nil {
		return responses.BadRequest(err)
	}
	data, err := app.Services.Budget.GetMonthly(context.Background(), userID, year, month)
	if err != nil {
		if errors.Is(err, budget.ErrInvalidPeriod()) {
			return responses.BadRequest(err)
		}
		return responses.InternalServerError(err)
	}
	return c.JSON(data)
}
