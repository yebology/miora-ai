package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterWatchlistProtectedRoutes registers watchlist routes (requires auth).
//
// Routes:
//
//	POST   /watchlist/follow    → follow a wallet
//	DELETE /watchlist/:address  → unfollow a wallet
//	GET    /watchlist           → list followed wallets
func RegisterWatchlistProtectedRoutes(r fiber.Router, h *handlers.WatchlistHandler) {

	wl := r.Group("/watchlist")
	wl.Post("/follow", h.Follow)
	wl.Put("/:address", h.Update)
	wl.Delete("/:address", h.Unfollow)
	wl.Get("/", h.List)

}
