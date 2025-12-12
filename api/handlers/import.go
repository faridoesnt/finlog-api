package handlers

import (
	"errors"
	"strconv"
	"time"

	"finlog-api/api/models/request"
	"finlog-api/api/models/responses"
	"finlog-api/api/services/importbatch"

	"github.com/gofiber/fiber/v2"
)

type importHistoryResponse struct {
	BatchID   int64     `json:"batch_id"`
	BatchSize int       `json:"batch_size"`
	CreatedAt time.Time `json:"created_at"`
}

// ImportTransactions stores encrypted import batches.
func ImportTransactions(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	var payload request.ImportBatchRequest
	if err := c.BodyParser(&payload); err != nil {
		return responses.BadRequest(err)
	}

	if err := app.Services.Import.StoreBatch(c.Context(), userID, payload); err != nil {
		if errors.Is(err, importbatch.ErrInvalidImportInput) {
			return responses.BadRequest(err)
		}
		if errors.Is(err, importbatch.ErrRateLimitExceeded) {
			return fiber.NewError(fiber.StatusTooManyRequests, err.Error())
		}
		return responses.InternalServerError(err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ImportHistory lists past import batches for the current user.
func ImportHistory(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	batches, err := app.Services.Import.ListHistory(c.Context(), userID)
	if err != nil {
		return responses.InternalServerError(err)
	}

	result := make([]importHistoryResponse, len(batches))
	for i, batch := range batches {
		result[i] = importHistoryResponse{
			BatchID:   batch.ID,
			BatchSize: batch.BatchSize,
			CreatedAt: batch.CreatedAt,
		}
	}
	return c.JSON(result)
}

// UndoImportBatch deletes all transactions from a specific batch.
func UndoImportBatch(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	idParam := c.Params("batch_id")
	batchID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return responses.BadRequest(err)
	}

	deleted, err := app.Services.Import.UndoBatch(c.Context(), userID, batchID)
	if err != nil {
		switch {
		case errors.Is(err, importbatch.ErrInvalidImportInput):
			return responses.BadRequest(err)
		case errors.Is(err, importbatch.ErrUndoRateLimitExceeded):
			return fiber.NewError(fiber.StatusTooManyRequests, err.Error())
		case errors.Is(err, importbatch.ErrImportBatchNotFound):
			return responses.NotFound(err)
		default:
			return responses.InternalServerError(err)
		}
	}

	return c.JSON(fiber.Map{
		"deleted_count": deleted,
		"batch_id":      batchID,
	})
}
