package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterAuthProtectedRoutes registers auth routes (requires wallet auth).
//
// Routes:
//
//	GET /auth/me       — get or create current user from wallet address
//	PUT /auth/email    — update email for notifications (optional)
func RegisterAuthProtectedRoutes(r fiber.Router, h *handlers.AuthHandler) {
	auth := r.Group("/auth")
	auth.Get("/me", h.Me)
	auth.Put("/email", h.UpdateEmail)
}
