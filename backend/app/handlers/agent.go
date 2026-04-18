// Package handlers contains the AI trading agent HTTP request handlers.
//
// All agent endpoints require wallet authentication.
// The user is identified by wallet address from the WalletAuth middleware.
package handlers

import (
	"strconv"

	"miora-ai/app/dto/requests"
	"miora-ai/app/interfaces"
	"miora-ai/app/output"
	"miora-ai/app/services"
	"miora-ai/constants"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

// AgentHandler handles agent-related HTTP requests.
type AgentHandler struct {
	agentService interfaces.IAgentService
	userService  *services.UserService
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(agentService interfaces.IAgentService, userService *services.UserService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		userService:  userService,
	}
}

// getUserID extracts the user ID from the wallet auth context.
func (h *AgentHandler) getUserID(c *fiber.Ctx) (uint, *fiber.Error) {
	walletAddress, _ := c.Locals("wallet_address").(string)
	if walletAddress == "" {
		return 0, fiber.NewError(fiber.StatusUnauthorized, constants.Unauthorized)
	}

	user, appErr := h.userService.FindOrCreateByWallet(walletAddress)
	if appErr != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, constants.Unauthorized)
	}

	return user.ID, nil
}

// GetStatus handles GET /agent/status.
func (h *AgentHandler) GetStatus(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	config, appErr := h.agentService.GetStatus(userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, config)
}

// UpdateConfig handles PUT /agent/config.
func (h *AgentHandler) UpdateConfig(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	var req requests.UpdateAgentConfig
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	config, appErr := h.agentService.UpdateConfig(
		userID, req.Budget, req.MaxPerTrade, req.RiskTolerance, req.MinScore, req.Conditions,
	)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Agent configuration updated.", config)
}

// Start handles POST /agent/start.
func (h *AgentHandler) Start(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	config, appErr := h.agentService.Start(userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Agent started.", config)
}

// Pause handles POST /agent/pause.
func (h *AgentHandler) Pause(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	config, appErr := h.agentService.Pause(userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Agent paused.", config)
}

// GetTrades handles GET /agent/trades.
func (h *AgentHandler) GetTrades(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	trades, appErr := h.agentService.GetTrades(userID, limit)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, trades)
}
