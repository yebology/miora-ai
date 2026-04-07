package interfaces

import "miora-ai/app/dto/responses"

// ISwapClient defines the contract for DEX aggregator quote fetching.
type ISwapClient interface {
	GetQuote(inputMint, outputMint, amount string, slippage int, chain ...string) (*responses.SwapQuote, error)
}
