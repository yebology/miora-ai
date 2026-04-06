package interfaces

import "miora-ai/app/dto"

// IDexScreener defines the contract for fetching token pair data.
type IDexScreener interface {
	GetTokenPairs(chainID, tokenAddress string) ([]dto.TokenPairData, error)
}
