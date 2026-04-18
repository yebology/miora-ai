// Package http registers agent-related routes (requires wallet auth).
package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterAgentProtectedRoutes registers AI trading agent routes (requires auth).
//
// Routes:
//
//	GET    /agent/bots           → list all bots
//	POST   /agent/bots           → create a new bot
//	GET    /agent/bots/:id       → get bot details
//	PUT    /agent/bots/:id       → update bot config
//	DELETE /agent/bots/:id       → delete a bot
//	POST   /agent/bots/:id/start → start a bot
//	POST   /agent/bots/:id/pause → pause a bot
//	GET    /agent/bots/:id/trades → get bot trade history
func RegisterAgentProtectedRoutes(r fiber.Router, h *handlers.AgentHandler) {
	bots := r.Group("/agent/bots")
	bots.Get("/", h.ListBots)
	bots.Post("/", h.CreateBot)
	bots.Get("/:id", h.GetBot)
	bots.Put("/:id", h.UpdateBot)
	bots.Delete("/:id", h.DeleteBot)
	bots.Post("/:id/start", h.StartBot)
	bots.Post("/:id/pause", h.PauseBot)
	bots.Get("/:id/trades", h.GetTrades)
}
