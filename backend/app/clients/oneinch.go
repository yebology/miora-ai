// 1inch client for EVM swap quotes.
// Uses 1inch Swap API (requires API key from portal.1inch.dev).
//
// API docs: https://portal.1inch.dev/documentation/apis/swap/classic-swap/quick-start
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/dto/responses"
	"miora-ai/constants"
)

// OneInch fetches swap quotes from the 1inch aggregator (EVM).
type OneInch struct {
	apiKey string
}

// NewOneInch creates a new 1inch client.
func NewOneInch(apiKey string) *OneInch {

	return &OneInch{apiKey: apiKey}

}

type oneInchQuoteResponse struct {
	DstAmount string `json:"dstAmount"`
}

// GetQuote fetches a swap quote from 1inch.
// chain parameter determines which EVM network to use.
func (o *OneInch) GetQuote(inputMint, outputMint, amount string, slippage int, chain ...string) (*responses.SwapQuote, error) {

	chainKey := "ethereum"
	if len(chain) > 0 && chain[0] != "" {
		chainKey = chain[0]
	}

	cfg := constants.GetChainConfig(chainKey)
	if cfg == nil || cfg.OneInchChainID == "" {
		return nil, fmt.Errorf("unsupported chain for 1inch: %s", chainKey)
	}

	url := fmt.Sprintf(
		"https://api.1inch.dev/swap/v6.0/%s/quote?src=%s&dst=%s&amount=%s",
		cfg.OneInchChainID, inputMint, outputMint, amount,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("1inch request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("1inch error %d: %s", resp.StatusCode, string(body))
	}

	var result oneInchQuoteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &responses.SwapQuote{
		Chain:        "evm",
		InputMint:    inputMint,
		OutputMint:   outputMint,
		InputAmount:  amount,
		OutputAmount: result.DstAmount,
	}, nil

}
