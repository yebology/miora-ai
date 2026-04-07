package handlers

import (
	"miora-ai/app/dto/requests"
	"miora-ai/app/entities"
	"miora-ai/app/output"
	"miora-ai/app/services"
	"miora-ai/constants"
	"miora-ai/pkg"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

type WatchlistHandler struct {
	service     *services.WatchlistService
	userService *services.UserService
}

func NewWatchlistHandler(service *services.WatchlistService, userService *services.UserService) *WatchlistHandler {

	return &WatchlistHandler{service: service, userService: userService}

}

// Follow handles POST /watchlist/follow.
func (h *WatchlistHandler) Follow(c *fiber.Ctx) error {

	user, appErr := h.getUser(c)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	var req requests.FollowWallet
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	if appErr := h.service.Follow(user.ID, req.WalletAddress, req.Chain, req.Recommendation, req.Conditions, req.EmailNotify); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessInsertData, nil)

}

// Unfollow handles DELETE /watchlist/:address.
func (h *WatchlistHandler) Unfollow(c *fiber.Ctx) error {

	user, appErr := h.getUser(c)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	address := c.Params("address")
	if address == "" {
		return output.GetError(c, fiber.StatusBadRequest, constants.AddressRequired)
	}

	if appErr := h.service.Unfollow(user.ID, address); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessDeleteData, nil)

}

// List handles GET /watchlist.
func (h *WatchlistHandler) List(c *fiber.Ctx) error {

	user, appErr := h.getUser(c)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	items, appErr := h.service.GetUserWatchlist(user.ID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetAllData, items)

}

// getUser extracts the current user from Firebase locals.
func (h *WatchlistHandler) getUser(c *fiber.Ctx) (*entities.User, *pkg.AppError) {

	uid, _ := c.Locals("firebase_uid").(string)
	if uid == "" {
		return nil, pkg.ErrUnauthorized(constants.Unauthorized)
	}

	user, appErr := h.userService.GetByFirebaseUID(uid)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil

}
