package handlers

import (
	"miora-ai/app/dto/requests"
	"miora-ai/app/services"
	"miora-ai/pkg"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	service *services.WalletService
}

func NewWalletHandler(service *services.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) Analyze(c *fiber.Ctx) error {
	var req requests.AnalyzeWallet
	if err := c.BodyParser(&req); err != nil {
		return pkg.ErrorResponse(c, pkg.NewAppError(fiber.StatusBadRequest, "invalid request body"))
	}

	if req.Address == "" || req.Chain == "" {
		return pkg.ErrorResponse(c, pkg.NewAppError(fiber.StatusBadRequest, "address and chain are required"))
	}

	if req.Chain != "evm" && req.Chain != "svm" {
		return pkg.ErrorResponse(c, pkg.NewAppError(fiber.StatusBadRequest, "chain must be 'evm' or 'svm'"))
	}

	result, err := h.service.AnalyzeWallet(req.Address, req.Chain)
	if err != nil {
		return pkg.ErrorResponse(c, pkg.NewAppError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(result)
}
