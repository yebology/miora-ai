// Package handlers contains the reputation HTTP request handler.
//
// Provides endpoints for querying on-chain trading reputation attestations
// published via EAS on Base Sepolia.
package handlers

import (
	"fmt"

	"miora-ai/app/dto/responses"
	"miora-ai/app/interfaces"
	"miora-ai/app/output"
	"miora-ai/constants"

	"github.com/gofiber/fiber/v2"
)

// ReputationHandler handles reputation-related HTTP requests.
type ReputationHandler struct {
	walletRepo interfaces.IWalletRepository
	easClient  interfaces.IEASClient
}

// NewReputationHandler creates a new ReputationHandler.
func NewReputationHandler(walletRepo interfaces.IWalletRepository, easClient interfaces.IEASClient) *ReputationHandler {
	return &ReputationHandler{
		walletRepo: walletRepo,
		easClient:  easClient,
	}
}

// explorerURL builds the EAS attestation explorer URL for a given attestation UID.
func explorerURL(uid string) string {
	return fmt.Sprintf("%s/attestation/view/%s", constants.BaseSepoliaEASScanURL, uid)
}

// GetReputation handles GET /reputation/:address.
// Returns the on-chain attestation data for a previously analyzed wallet.
func (h *ReputationHandler) GetReputation(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return output.GetError(c, fiber.StatusBadRequest, constants.AddressRequired)
	}

	// Look up wallet in database
	wallet, err := h.walletRepo.FindByAddress(address)
	if err != nil || wallet == nil {
		return output.GetError(c, fiber.StatusNotFound, "Wallet not found. Analyze it first via POST /api/wallets/analyze.")
	}

	// Get stored metric with attestation UID
	metric, err := h.walletRepo.GetMetric(wallet.ID)
	if err != nil || metric == nil {
		return output.GetError(c, fiber.StatusNotFound, "No scoring data found for this wallet.")
	}

	if metric.AttestationUID == "" {
		return output.GetError(c, fiber.StatusNotFound, "No on-chain attestation found for this wallet. Re-analyze to publish attestation.")
	}

	// Query on-chain attestation
	attestation, err := h.easClient.GetAttestation(metric.AttestationUID)
	if err != nil {
		// Fall back to database data if on-chain query fails
		return output.GetSuccess(c, constants.SuccessGetData, responses.Reputation{
			Address:           address,
			Chain:             wallet.Chain,
			Score:             uint8(metric.FinalScore),
			Recommendation:    metric.Recommendation,
			TotalTransactions: uint32(metric.TotalTransactions),
			AttestationUID:    metric.AttestationUID,
			AttestationTxHash: metric.AttestationTxHash,
			ExplorerURL:       explorerURL(metric.AttestationUID),
		})
	}

	return output.GetSuccess(c, constants.SuccessGetData, responses.Reputation{
		Address:           address,
		Chain:             attestation.Chain,
		Score:             attestation.Score,
		Recommendation:    attestation.Recommendation,
		TotalTransactions: attestation.TotalTransactions,
		AttestationUID:    attestation.UID,
		AttestationTxHash: metric.AttestationTxHash,
		Attester:          attestation.Attester,
		Timestamp:         attestation.Time,
		ExplorerURL:       explorerURL(attestation.UID),
	})
}
