package responses

// SwapQuote is the response from POST /api/swap/quote.
//
// Example response:
//
//	{
//	  "chain": "base",
//	  "input_mint": "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
//	  "output_mint": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
//	  "input_amount": "1000000000000000000",
//	  "output_amount": "1847500000",
//	  "price_impact": "0.01",
//	  "route": "Uniswap V3 → SushiSwap"
//	}
type SwapQuote struct {
	Chain        string `json:"chain"`                  // EVM chain name
	InputMint    string `json:"input_mint"`             // Token address being sold
	OutputMint   string `json:"output_mint"`            // Token address being bought
	InputAmount  string `json:"input_amount"`           // Amount being sold (smallest unit)
	OutputAmount string `json:"output_amount"`          // Amount received (smallest unit)
	PriceImpact  string `json:"price_impact,omitempty"` // How much the swap moves the market price (e.g. "0.01" = 0.01%). Higher = worse.
	Route        string `json:"route,omitempty"`        // DEX route taken (e.g. "Raydium → Orca"). Shows which DEXes the aggregator used.
}
