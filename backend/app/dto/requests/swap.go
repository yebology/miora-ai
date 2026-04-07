package requests

// SwapQuote is the request body for POST /api/swap/quote.
//
// Example request:
//
//	{
//	  "chain": "svm",
//	  "input_mint": "So11111111111111111111111111111111111111112",
//	  "output_mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
//	  "amount": "100000000",
//	  "slippage": 50
//	}
type SwapQuote struct {
	Chain      string `json:"chain" validate:"required,oneof=evm svm ethereum arbitrum optimism base polygon solana"` // Blockchain: "evm" (Ethereum) or "svm" (Solana)
	InputMint  string `json:"input_mint" validate:"required"`                                                         // Token address to sell (e.g. SOL mint address, or ETH contract address)
	OutputMint string `json:"output_mint" validate:"required"`                                                        // Token address to buy
	Amount     string `json:"amount" validate:"required"`                                                             // Amount to swap in smallest unit (lamports for SOL, wei for ETH)
	Slippage   int    `json:"slippage"`                                                                               // Max acceptable price difference in basis points (50 = 0.5%, 100 = 1%). Optional, defaults to 50.
}
