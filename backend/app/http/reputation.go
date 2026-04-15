// Package http registers reputation-related routes.
package http

import (
	"miora-ai/app/handlers"
	"miora-ai/app/middleware"

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

// RegisterReputationX402Routes registers the x402-protected reputation query route.
//
// Routes:
//
//	GET /reputation/query?address=0x... — query reputation score (requires x402 USDC payment)
func RegisterReputationX402Routes(r fiber.Router, h *handlers.ReputationHandler, x402Cfg middleware.X402Config) {
	reputation := r.Group("/reputation")
	reputation.Get("/query", middleware.X402Middleware(x402Cfg), h.QueryReputation)
}
