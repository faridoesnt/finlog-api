package handlers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"finlog-api/api/models/request"
	"finlog-api/api/models/responses"
	"finlog-api/api/services/transaction"
)

// GetRecentTransactions returns the latest transactions for a given month.
func GetRecentTransactions(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	year, month, err := parsePeriod(c)
	if err != nil {
		return responses.BadRequest(err)
	}
	txs, err := app.Services.Transactions.GetRecentTransactions(context.Background(), userID, year, month)
	if err != nil {
		return responses.InternalServerError(err)
	}
	return c.JSON(txs)
}

// GetTransactions returns transactions for the given period.
func GetTransactions(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	year, month, err := parsePeriod(c)
	if err != nil {
		return responses.BadRequest(err)
	}
	txs, err := app.Services.Transactions.GetTransactions(context.Background(), userID, year, month)
	if err != nil {
		return responses.InternalServerError(err)
	}
	return c.JSON(txs)
}

// CreateTransaction stores a new transaction.
func CreateTransaction(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)

	body := request.CreateTransaction{}
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}

	occurredAt, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		return responses.BadRequest(errors.New("invalid date format"))
	}

	body.OccurredAt = occurredAt

	tx, err := app.Services.Transactions.CreateTransaction(context.Background(), userID, body)
	if err != nil {
		return mapTransactionError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}

// UpdateTransaction updates a transaction by id.
func UpdateTransaction(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return responses.BadRequest(errors.New("invalid transaction id"))
	}

	body := request.CreateTransaction{}
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}

	occurredAt, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		return responses.BadRequest(errors.New("invalid date format"))
	}

	body.OccurredAt = occurredAt

	if err := app.Services.Transactions.UpdateTransaction(context.Background(), userID, id, body); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

// UpdateTransactionNotes updates notes for multiple transactions.
func UpdateTransactionNotes(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	type req struct {
		TransactionIDs []string `json:"transactionIds"`
		Notes          string   `json:"notes"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	ids, err := parseIDs(body.TransactionIDs)
	if err != nil {
		return responses.BadRequest(err)
	}
	if err := app.Services.Transactions.UpdateNotes(context.Background(), userID, ids, body.Notes); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

// UpdateTransactionAmount updates amount for multiple transactions.
func UpdateTransactionAmount(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	type req struct {
		TransactionIDs []string `json:"transactionIds"`
		Amount         int64    `json:"amount"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	ids, err := parseIDs(body.TransactionIDs)
	if err != nil {
		return responses.BadRequest(err)
	}
	if err := app.Services.Transactions.UpdateAmount(context.Background(), userID, ids, body.Amount); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

// UpdateTransactionDate updates occurrence date for multiple transactions.
func UpdateTransactionDate(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	type req struct {
		TransactionIDs []string `json:"transactionIds"`
		Date           string   `json:"date"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	ids, err := parseIDs(body.TransactionIDs)
	if err != nil {
		return responses.BadRequest(err)
	}
	date, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		return responses.BadRequest(errors.New("invalid date format"))
	}
	if err := app.Services.Transactions.UpdateDate(context.Background(), userID, ids, date); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteTransaction deletes a transaction by id.
func DeleteTransaction(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return responses.BadRequest(errors.New("invalid transaction id"))
	}
	if err := app.Services.Transactions.DeleteTransaction(context.Background(), userID, id); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteTransactions deletes transactions in bulk.
func DeleteTransactions(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	type req struct {
		TransactionIDs []string `json:"transactionIds"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return responses.BadRequest(err)
	}
	ids, err := parseIDs(body.TransactionIDs)
	if err != nil {
		return responses.BadRequest(err)
	}
	if err := app.Services.Transactions.DeleteTransactions(context.Background(), userID, ids); err != nil {
		return mapTransactionError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parsePeriod(c *fiber.Ctx) (int, int, error) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return 0, 0, errors.New("invalid year")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return 0, 0, errors.New("invalid month")
	}
	return year, month, nil
}

func parseIDs(ids []string) ([]int64, error) {
	out := make([]int64, 0, len(ids))
	for _, raw := range ids {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, errors.New("invalid id in list")
		}
		out = append(out, id)
	}
	return out, nil
}

func mapTransactionError(err error) error {
	switch {
	case errors.Is(err, transaction.ErrInvalidTransaction()):
		return responses.BadRequest(err)
	case errors.Is(err, transaction.ErrTransactionNotFound()):
		return responses.NotFound(err)
	case errors.Is(err, transaction.ErrCategoryNotFound()):
		return responses.BadRequest(err)
	case errors.Is(err, transaction.ErrUnsupportedBulk()):
		return responses.BadRequest(err)
	default:
		return responses.BadRequest(err)
	}
}
