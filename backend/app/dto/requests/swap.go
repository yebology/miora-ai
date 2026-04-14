package requests

// SwapQuote is the request body for POST /api/swap/quote.
//
// Example request:
//
//	{
//	  "chain": "base",
//	  "input_mint": "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
//	  "output_mint": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
//	  "amount": "1000000000000000000",
//	  "slippage": 50
//	}
type SwapQuote struct {
	Chain      string `json:"chain" validate:"required,oneof=evm ethereum arbitrum optimism base polygon"` // EVM chain name
	InputMint  string `json:"input_mint" validate:"required"`                                              // Token address to sell
	OutputMint string `json:"output_mint" validate:"required"`                                             // Token address to buy
	Amount     string `json:"amount" validate:"required"`                                                  // Amount to swap in smallest unit (wei for ETH)
	Slippage   int    `json:"slippage"`                                                                    // Max acceptable price difference in basis points (50 = 0.5%, 100 = 1%). Optional, defaults to 50.
}
