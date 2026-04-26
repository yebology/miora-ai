package interfaces

import "miora-ai/app/dto"

// IMoralis defines the contract for fetching token price data.
type IMoralis interface {
	// GetTokenPrice fetches the USD price of a token at a specific block.
	// If block is 0, returns the current price.
	GetTokenPrice(chain, tokenAddress string, block uint64) (*dto.TokenPriceData, error)
}
