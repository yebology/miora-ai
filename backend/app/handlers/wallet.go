// Package handlers contains HTTP request handlers for each domain.
//
// Every handler follows this pattern:
//  1. Parse + validate request body via utils.ParseAndValidateBody()
//  2. Extract URL params or user context if needed
//  3. Call the service layer
//  4. Return output.GetSuccess() or output.GetError()
//
// Handlers never contain business logic — only parse, delegate, and respond.
// Each handler depends on a service interface (e.g. IWalletService), not a concrete struct.
package handlers

import (
	"miora-ai/app/dto/requests"
	"miora-ai/app/interfaces"
	"miora-ai/app/output"
	"miora-ai/constants"
	"miora-ai/utils"

	"github.com/gofiber/fiber/v2"
)

// WalletHandler handles wallet-related HTTP requests.
// Depends on IWalletService interface — no direct service struct dependency.
type WalletHandler struct {
	service interfaces.IWalletService
}

// NewWalletHandler creates a new WalletHandler with the given service.
func NewWalletHandler(service interfaces.IWalletService) *WalletHandler {

	return &WalletHandler{service: service}

}

// Analyze handles POST /wallets/analyze.
//
// Request body:
//
//	{
//	  "address": "0x...",   // required — wallet address to analyze
//	  "chain":   "evm"      // required — must be "evm" or "svm"
//	}
//
// Success response (200):
//
//	{
//	  "status": "success",
//	  "message": "Data retrieved successfully.",
//	  "data": { address, chain, total_transactions, scores..., recommendation }
//	}
//
// Error responses:
//   - 400: invalid request body or validation failed
//   - 502: failed to fetch data from Alchemy
//   - 500: internal server error
func (h *WalletHandler) Analyze(c *fiber.Ctx) error {

	// 1. Parse + validate request body
	var req requests.AnalyzeWallet
	if appErr := utils.ParseAndValidateBody(c, &req); appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	// 2. Call service
	result, appErr := h.service.AnalyzeWallet(req.Address, req.Chain, req.Limit)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	// 3. Return success
	return output.GetSuccess(c, constants.SuccessGetData, result)

}

// GetWallet handles GET /wallets/:address.
// Retrieves a previously analyzed wallet from the database.
func (h *WalletHandler) GetWallet(c *fiber.Ctx) error {

	address := c.Params("address")
	if address == "" {
		return output.GetError(c, fiber.StatusBadRequest, constants.AddressRequired)
	}

	result, appErr := h.service.GetWallet(address)
	if appErr != nil {
		return output.GetError(c, appErr.Code, appErr.Message)
	}

	return output.GetSuccess(c, constants.SuccessGetData, result)

}
