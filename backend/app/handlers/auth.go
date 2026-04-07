package handlers

import (
	"miora-ai/app/output"
	"miora-ai/app/services"
	"miora-ai/constants"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {

	return &AuthHandler{service: service}

}

// Me handles GET /auth/me.
// Returns the current user based on Firebase token (set by middleware).
// Creates the user in DB if first login.
func (h *AuthHandler) Me(c *fiber.Ctx) error {

	uid, _ := c.Locals("firebase_uid").(string)
	email, _ := c.Locals("firebase_email").(string)
	name, _ := c.Locals("firebase_name").(string)
	avatar, _ := c.Locals("firebase_avatar").(string)

	if uid == "" {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	user, appErr := h.service.FindOrCreateFromFirebase(uid, email, name, avatar)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, fiber.Map{
		"id":     user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"avatar": user.Avatar,
	})

}
