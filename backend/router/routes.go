package router

import (
	apphttp "miora-ai/app/http"
	"miora-ai/app/middleware"
	"miora-ai/app/ws"
	"miora-ai/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetUp initializes the DI container and registers all API routes.
//
// Route structure:
//
//	/api
//	  GET  /health              → health check (public)
//	  POST /wallets/analyze     → analyze wallet (public)
//	  GET  /wallets/:address    → get stored analysis (public)
//	  GET  /reputation/:address → get on-chain reputation (public)
//	  GET  /auth/me             → get/create user (protected, wallet)
func SetUp(app *fiber.App, db *gorm.DB, cfg *config.Config, hub *ws.Hub) {

	api := app.Group("/api")
	container := NewContainer(db, cfg.AlchemyAPIKey, cfg.MoralisAPIKey, cfg.GeminiAPIKey, cfg.Scoring, cfg.EAS, hub)

	// Start wallet monitor in background
	go container.Monitor.Start()

	// Start agent trading loop in background
	go container.AgentLoop.Start()

	// Health
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// WebSocket
	app.Use("/ws", ws.UpgradeHandler())
	app.Get("/ws", ws.ConnectHandler(hub))

	// Public routes
	apphttp.RegisterWalletPublicRoutes(api, container.WalletHandler)
	apphttp.RegisterReputationPublicRoutes(api, container.ReputationHandler)

	// Protected routes (wallet auth — X-Wallet-Address header)
	protected := api.Group("", middleware.WalletAuth())
	apphttp.RegisterAuthProtectedRoutes(protected, container.AuthHandler)
	apphttp.RegisterWatchlistProtectedRoutes(protected, container.WatchlistHandler)
	apphttp.RegisterAgentProtectedRoutes(protected, container.AgentHandler)

}
