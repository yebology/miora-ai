package http

import (
	"miora-ai/app/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterAuthProtectedRoutes registers auth routes that require Firebase token.
//
// Routes:
//
//	GET /auth/me — get or create current user from Firebase token
func RegisterAuthProtectedRoutes(r fiber.Router, h *handlers.AuthHandler) {

	auth := r.Group("/auth")
	auth.Get("/me", h.Me)

}
