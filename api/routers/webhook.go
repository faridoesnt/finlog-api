package routers

import (
	"finlog-api/api/contracts"
	"finlog-api/api/handlers"
)

func RegisterWebhook(app *contracts.App) {
	app.Fiber.Post("/webhooks/resend", handlers.ResendWebhook(app))
}