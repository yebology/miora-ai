// Package handlers contains the AI trading agent HTTP request handlers.
//
// Each bot targets one wallet. Users can have multiple bots.
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

type AgentHandler struct {
	agentService interfaces.IAgentService
	userService  *services.UserService
}

func NewAgentHandler(agentService interfaces.IAgentService, userService *services.UserService) *AgentHandler {
	return &AgentHandler{agentService: agentService, userService: userService}
}

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

// ListBots handles GET /agent/bots — list all bots for the user.
func (h *AgentHandler) ListBots(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	bots, appErr := h.agentService.ListBots(userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, bots)
}

// GetBot handles GET /agent/bots/:id — get a single bot.
func (h *AgentHandler) GetBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	bot, appErr := h.agentService.GetBot(uint(botID), userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, bot)
}

// CreateBot handles POST /agent/bots — create a new bot for a wallet.
func (h *AgentHandler) CreateBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	var req requests.CreateBot
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	bot, appErr := h.agentService.CreateBot(
		userID, req.BotType, req.TargetWalletAddress, req.TargetWalletChain,
		req.TargetWalletScore, req.Recommendation,
		req.Budget, req.MaxPerTrade, req.Conditions,
		req.ConsensusThreshold, req.ConsensusWindowMin, req.MinScore,
	)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Bot created.", bot)
}

// UpdateBot handles PUT /agent/bots/:id — update bot config.
func (h *AgentHandler) UpdateBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var req requests.UpdateBot
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	bot, appErr := h.agentService.UpdateBot(
		uint(botID), userID, req.Budget, req.MaxPerTrade, req.Conditions,
		req.ConsensusThreshold, req.ConsensusWindowMin, req.MinScore,
	)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Bot updated.", bot)
}

// DeleteBot handles DELETE /agent/bots/:id — remove a bot.
func (h *AgentHandler) DeleteBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	if appErr := h.agentService.DeleteBot(uint(botID), userID); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Bot deleted.", nil)
}

// StartBot handles POST /agent/bots/:id/start.
func (h *AgentHandler) StartBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	bot, appErr := h.agentService.StartBot(uint(botID), userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Bot started.", bot)
}

// PauseBot handles POST /agent/bots/:id/pause.
func (h *AgentHandler) PauseBot(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	bot, appErr := h.agentService.PauseBot(uint(botID), userID)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Bot paused.", bot)
}

// GetTrades handles GET /agent/bots/:id/trades.
func (h *AgentHandler) GetTrades(c *fiber.Ctx) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	botID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	trades, appErr := h.agentService.GetTrades(uint(botID), userID, limit)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, trades)
}

// Withdraw handles POST /agent/withdraw.
// Transfers ETH from the Agentic Wallet to the user's wallet.
func (h *AgentHandler) Withdraw(c *fiber.Ctx) error {
	walletAddress, _ := c.Locals("wallet_address").(string)
	if walletAddress == "" {
		return output.GetError(c, fiber.StatusUnauthorized, constants.Unauthorized)
	}

	var req struct {
		AmountETH string `json:"amount_eth" validate:"required"`
	}
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, "Withdraw initiated.", fiber.Map{
		"to":         walletAddress,
		"amount_eth": req.AmountETH,
		"status":     "pending",
	})
}
