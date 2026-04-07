// Jupiter client for Solana swap quotes.
// Uses Jupiter Swap API v1 (free, no API key required).
//
// API docs: https://dev.jup.ag/docs/swap/get-quote
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/dto/responses"
)

// Jupiter fetches swap quotes from the Jupiter aggregator (Solana).
type Jupiter struct{}

// NewJupiter creates a new Jupiter client.
func NewJupiter() *Jupiter {

	return &Jupiter{}

}

type jupiterQuoteResponse struct {
	InputMint      string `json:"inputMint"`
	OutputMint     string `json:"outputMint"`
	InAmount       string `json:"inAmount"`
	OutAmount      string `json:"outAmount"`
	PriceImpactPct string `json:"priceImpactPct"`
	RoutePlan      []struct {
		SwapInfo struct {
			Label string `json:"label"`
		} `json:"swapInfo"`
	} `json:"routePlan"`
}

// GetQuote fetches a swap quote from Jupiter.
// amount is in smallest unit (lamports for SOL, token decimals for SPL tokens).
// slippage is in basis points (50 = 0.5%).
//
// Endpoint: GET https://api.jup.ag/swap/v1/quote
func (j *Jupiter) GetQuote(inputMint, outputMint, amount string, slippage int, chain ...string) (*responses.SwapQuote, error) {

	if slippage == 0 {
		slippage = 50
	}

	url := fmt.Sprintf(
		"https://api.jup.ag/swap/v1/quote?inputMint=%s&outputMint=%s&amount=%s&slippageBps=%d",
		inputMint, outputMint, amount, slippage,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("jupiter request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("jupiter error %d: %s", resp.StatusCode, string(body))
	}

	var result jupiterQuoteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Extract route labels
	route := ""
	for i, r := range result.RoutePlan {
		if i > 0 {
			route += " → "
		}
		route += r.SwapInfo.Label
	}

	return &responses.SwapQuote{
		Chain:        "svm",
		InputMint:    result.InputMint,
		OutputMint:   result.OutputMint,
		InputAmount:  result.InAmount,
		OutputAmount: result.OutAmount,
		PriceImpact:  result.PriceImpactPct,
		Route:        route,
	}, nil

}
