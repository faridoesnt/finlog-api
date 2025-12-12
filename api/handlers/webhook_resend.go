// handlers/resend_webhook.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/models/request"

	"github.com/gofiber/fiber/v2"
	svix "github.com/svix/svix-webhooks/go"
)

func ResendWebhook(app *contracts.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := app.Config[constants.ResendWebhookSecret]
		if secret == "" {
			return c.Status(http.StatusInternalServerError).SendString("webhook secret not configured")
		}

		payload := c.Body()

		svixID := c.Get("svix-id")
		svixTS := c.Get("svix-timestamp")
		svixSig := c.Get("svix-signature")
		if svixID == "" || svixTS == "" || svixSig == "" {
			return c.Status(http.StatusBadRequest).SendString("missing svix headers")
		}

		wh, err := svix.NewWebhook(secret)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("webhook init failed")
		}

		headers := http.Header{}
		headers.Add("svix-id", svixID)
		headers.Add("svix-timestamp", svixTS)
		headers.Add("svix-signature", svixSig)

		if err := wh.Verify(payload, headers); err != nil {
			app.Logger.Warn().Err(err).Msg("resend_webhook_invalid_signature")
			return c.Status(http.StatusBadRequest).SendString("invalid signature")
		}

		var body request.ResendWebhookPayload
		if err := json.Unmarshal(payload, &body); err != nil {
			return c.Status(http.StatusBadRequest).SendString("invalid json")
		}

		app.Logger.Info().
			Str("type", body.Type).
			Str("email_id", body.Data.EmailID).
			Any("to", body.Data.To).
			Msg("resend_webhook_verified")

		if body.Type == "" || body.Data.EmailID == "" {
			return c.Status(http.StatusBadRequest).SendString("invalid payload")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.Services.Email.HandleWebhook(ctx, body, payload); err != nil {
			app.Logger.Error().
				Err(err).
				Str("type", body.Type).
				Str("email_id", body.Data.EmailID).
				Msg("resend_webhook_handle_failed")
			// webhook should be 2xx if you want Resend to not retry,
			// but for now it's safe: 500 so Resend can retry.
			return c.Status(http.StatusInternalServerError).SendString("failed")
		}

		return c.SendStatus(http.StatusOK)
	}
}
