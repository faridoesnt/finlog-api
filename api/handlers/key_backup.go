package handlers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"

	"finlog-api/api/models/request"
	"finlog-api/api/models/responses"
	"finlog-api/api/services/keybackup"
)

// StoreKeyBackup persists an encrypted key backup when the client provisions a vault for the first time.
func StoreKeyBackup(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	payload := request.DataKeyBackupPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return responses.BadRequest(err)
	}

	key, err := app.Services.KeyBackup.StoreKeyBackup(context.Background(), userID, payload.EncryptedDataKey, payload.Salt)
	if err != nil {
		return mapKeyBackupError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         key.ID,
		"salt":       key.Salt,
		"created_at": key.CreatedAt,
	})
}

// RotateKeyBackup replaces an active key with a new encrypted backup while keeping rotation history.
func RotateKeyBackup(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	payload := request.DataKeyBackupPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return responses.BadRequest(err)
	}

	key, err := app.Services.KeyBackup.RotateKey(context.Background(), userID, payload.EncryptedDataKey, payload.Salt)
	if err != nil {
		return mapKeyBackupError(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         key.ID,
		"salt":       key.Salt,
		"created_at": key.CreatedAt,
	})
}

// GetActiveKeyBackup returns the active encrypted data key along with the associated salt for recovery.
func GetActiveKeyBackup(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	key, err := app.Services.KeyBackup.GetActiveKey(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.NotFound(err)
		}
		return responses.InternalServerError(err)
	}
	return c.JSON(fiber.Map{
		"encrypted_data_key": key.EncryptedDataKey,
		"salt":               key.Salt,
		"created_at":         key.CreatedAt,
		"rotated_at":         key.RotatedAt,
	})
}

// GetKeyBackupStatus exposes audit-friendly metadata about the user's key material.
func GetKeyBackupStatus(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int64)
	status, err := app.Services.KeyBackup.GetKeyStatus(context.Background(), userID)
	if err != nil {
		return responses.InternalServerError(err)
	}
	return c.JSON(status)
}

func mapKeyBackupError(err error) error {
	switch {
	case errors.Is(err, keybackup.ErrInvalidPayload()):
		return responses.BadRequest(err)
	case errors.Is(err, keybackup.ErrActiveKeyExists()):
		return responses.Conflict(err)
	case errors.Is(err, keybackup.ErrKeyNotFound()):
		return responses.NotFound(err)
	default:
		return responses.InternalServerError(err)
	}
}
