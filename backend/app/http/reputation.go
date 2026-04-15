// Package http registers reputation-related routes.
package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterReputationPublicRoutes registers reputation routes (public, no auth required).
//
// Routes:
//
//	GET /reputation/:address — get on-chain reputation attestation for a wallet
func RegisterReputationPublicRoutes(r fiber.Router, h *handlers.ReputationHandler) {
	reputation := r.Group("/reputation")
	reputation.Get("/:address", h.GetReputation)
}
