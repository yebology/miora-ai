package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterSwapPublicRoutes registers swap-related routes.
//
// Routes:
//
//	POST /swap/quote — get a swap quote from 1inch (EVM)
func RegisterSwapPublicRoutes(r fiber.Router, h *handlers.SwapHandler) {

	swap := r.Group("/swap")
	swap.Post("/quote", h.Quote)

}
