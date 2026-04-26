// Package handlers contains the auth HTTP request handler.
//
// Auth is wallet-based — user connects MetaMask, frontend sends wallet address.
package handlers

import (
	"miora-ai/app/output"
	"miora-ai/app/services"
	"miora-ai/constants"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Me handles GET /auth/me.
// Returns the current user based on wallet address (set by WalletAuth middleware).
// Creates the user in DB if first time.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	walletAddress, _ := c.Locals("wallet_address").(string)
	if walletAddress == "" {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	user, appErr := h.service.FindOrCreateByWallet(walletAddress)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, user)
}

// UpdateEmail handles PUT /auth/email.
// Sets the user's email for trade notifications (optional).
func (h *AuthHandler) UpdateEmail(c *fiber.Ctx) error {
	walletAddress, _ := c.Locals("wallet_address").(string)
	if walletAddress == "" {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	user, appErr := h.service.UpdateEmail(walletAddress, req.Email)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Email updated.", user)
}
