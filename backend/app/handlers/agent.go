// Package handlers contains the AI trading agent HTTP request handlers.
//
// All agent endpoints require Firebase authentication.
// The user ID is extracted from the Firebase middleware context.
package handlers

import (
	"strconv"

	"miora-ai/app/dto/requests"
	"miora-ai/app/interfaces"
	"miora-ai/app/output"
	"miora-ai/constants"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

// AgentHandler handles agent-related HTTP requests.
type AgentHandler struct {
	agentService interfaces.IAgentService
	userService  interfaces.IUserService
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(agentService interfaces.IAgentService, userService interfaces.IUserService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		userService:  userService,
	}
}

// getUserID extracts the user ID from the Firebase auth context.
func (h *AgentHandler) getUserID(c *fiber.Ctx) (uint, *fiber.Error) {
	firebaseUID, ok := c.Locals("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		return 0, fiber.NewError(fiber.StatusUnauthorized, constants.Unauthorized)
	}

	user, appErr := h.userService.GetByFirebaseUID(firebaseUID)
	if appErr != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, constants.Unauthorized)
	}

	return user.ID, nil
}

// GetStatus handles GET /agent/status.
// Returns the current agent configuration and status.
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
// Updates the agent configuration (budget, risk tolerance, conditions, etc.).
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
// Activates the AI trading agent.
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
// Pauses the AI trading agent.
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
// Returns the agent's trade history.
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
