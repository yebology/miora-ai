package handlers

import (
	"miora-ai/app/dto/requests"
	"miora-ai/app/output"
	"miora-ai/app/services"
	"miora-ai/constants"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

// SwapHandler handles swap-related HTTP requests.
type SwapHandler struct {
	service *services.SwapService
}

// NewSwapHandler creates a new SwapHandler.
func NewSwapHandler(service *services.SwapService) *SwapHandler {

	return &SwapHandler{service: service}

}

// Quote handles POST /swap/quote.
// Returns a swap quote without executing the swap.
func (h *SwapHandler) Quote(c *fiber.Ctx) error {

	var req requests.SwapQuote
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	result, appErr := h.service.GetQuote(req.Chain, req.InputMint, req.OutputMint, req.Amount, req.Slippage)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, result)

}
