// Package http registers agent-related routes (requires Firebase auth).
package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterAgentProtectedRoutes registers AI trading agent routes (requires auth).
//
// Routes:
//
//	GET    /agent/status  → get agent status and config
//	PUT    /agent/config  → update agent configuration
//	POST   /agent/start   → start the agent
//	POST   /agent/pause   → pause the agent
//	GET    /agent/trades  → get agent trade history
func RegisterAgentProtectedRoutes(r fiber.Router, h *handlers.AgentHandler) {
	agent := r.Group("/agent")
	agent.Get("/status", h.GetStatus)
	agent.Put("/config", h.UpdateConfig)
	agent.Post("/start", h.Start)
	agent.Post("/pause", h.Pause)
	agent.Get("/trades", h.GetTrades)
}
