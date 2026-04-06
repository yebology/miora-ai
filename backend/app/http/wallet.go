// Package http registers route groups for each domain.
//
// Each file in this package follows the pattern:
//
//	func RegisterXxxRoutes(r fiber.Router, h *handlers.XxxHandler) {
//	    group := r.Group("/xxx")
//	    group.Post("/", h.Create)
//	    group.Get("/", h.List)
//	}
//
// Route registration is called from router/routes.go during app setup.
package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterWalletPublicRoutes registers wallet-related routes that don't require authentication.
//
// Routes:
//
//	POST /wallets/analyze — analyze a wallet address and return scoring + recommendation
func RegisterWalletPublicRoutes(r fiber.Router, h *handlers.WalletHandler) {

	wallets := r.Group("/wallets")
	wallets.Post("/analyze", h.Analyze)
	wallets.Get("/:address", h.GetWallet)

}
