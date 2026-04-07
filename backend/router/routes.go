package router

import (
	apphttp "miora-ai/app/http"
	"miora-ai/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetUp initializes the DI container and registers all API routes.
//
// Route structure:
//
//	/api
//	  GET  /health              → health check
//	  POST /wallets/analyze     → analyze a wallet address
//
// To add a new route group:
//  1. Create app/http/xxx.go with RegisterXxxRoutes(r fiber.Router, h *handler)
//  2. Call it here with the appropriate handler from container
func SetUp(app *fiber.App, db *gorm.DB, alchemyAPIKey, moralisAPIKey, birdeyeAPIKey, geminiAPIKey, oneInchAPIKey string, scoring config.ScoringConfig) {

	api := app.Group("/api")
	container := NewContainer(db, alchemyAPIKey, moralisAPIKey, birdeyeAPIKey, geminiAPIKey, oneInchAPIKey, scoring)

	// Health
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public routes
	apphttp.RegisterWalletPublicRoutes(api, container.WalletHandler)
	apphttp.RegisterSwapPublicRoutes(api, container.SwapHandler)

}
