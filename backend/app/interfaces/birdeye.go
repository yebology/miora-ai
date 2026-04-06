package interfaces

import "miora-ai/app/dto"

// IBirdeye defines the contract for fetching Solana token historical prices.
// Used for PnL calculation on Solana where Moralis doesn't support historical prices.
type IBirdeye interface {
	// GetHistoricalPrice fetches the token price closest to the given unix timestamp.
	// Returns the price at that point in time.
	GetHistoricalPrice(tokenAddress string, unixTimestamp int64) (*dto.TokenPriceData, error)
}
