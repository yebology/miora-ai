package routes

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, walletHandler *handlers.WalletHandler) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	wallet := api.Group("/wallets")
	wallet.Post("/analyze", walletHandler.Analyze)
}
