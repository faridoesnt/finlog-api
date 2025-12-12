package routers

import (
	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/handlers"
	"finlog-api/api/middlewares"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Init(app *contracts.App) {
	app.Fiber.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Support both /api and /api/v1 prefixes for compatibility.
	registerAPIRoutes(app, app.Fiber.Group("/api"))
	registerAPIRoutes(app, app.Fiber.Group("/api/v1"))
	registerAPIRoutes(app, app.Fiber.Group("/v1"))
	RegisterWebhook(app)
}

func registerAPIRoutes(app *contracts.App, api fiber.Router) {
	authGroup := api.Group("/auth")
	authGroup.Post("/login", handlers.AuthLogin)
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/resend-verification", handlers.ResendVerification)
	authGroup.Get("/verify", handlers.VerifyEmail)
	authGroup.Post("/refresh", handlers.Refresh)

	jwtTTL := parseDuration(app.Config[constants.JWT_TTL], time.Hour)
	protected := api.Group("", middlewares.JWT([]byte(app.Config[constants.JWT_SECRET]), jwtTTL))

	protected.Post("/auth/logout", handlers.Logout)

	protected.Get("/categories", handlers.GetCategories)
	protected.Post("/categories", handlers.CreateCategory)
	protected.Put("/categories/:id", handlers.UpdateCategory)
	protected.Delete("/categories/:id", handlers.DeleteCategory)

	protected.Get("/recent-transactions", handlers.GetRecentTransactions)
	protected.Get("/transactions", handlers.GetTransactions)
	protected.Post("/transactions", handlers.CreateTransaction)
	protected.Post("/transactions/import", handlers.ImportTransactions)
	protected.Get("/transactions/import/history", handlers.ImportHistory)
	protected.Delete("/transactions/import/:batch_id", handlers.UndoImportBatch)
	protected.Put("/transactions/:id", handlers.UpdateTransaction)
	protected.Put("/transactions/bulk/notes", handlers.UpdateTransactionNotes)
	protected.Put("/transactions/bulk/amounts", handlers.UpdateTransactionAmount)
	protected.Put("/transactions/bulk/dates", handlers.UpdateTransactionDate)
	protected.Delete("/transactions/:id", handlers.DeleteTransaction)
	protected.Delete("/transactions/bulk/delete", handlers.DeleteTransactions)

	protected.Get("/budget", handlers.GetBudget)

	keyGroup := protected.Group("/keys")
	keyGroup.Post("/backup", handlers.StoreKeyBackup)
	keyGroup.Put("/backup/rotate", handlers.RotateKeyBackup)
	keyGroup.Get("/backup", handlers.GetActiveKeyBackup)
	keyGroup.Get("/backup/status", handlers.GetKeyBackupStatus)
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	return fallback
}
