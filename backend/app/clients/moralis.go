// Moralis client for fetching token prices.
//
// Supports EVM chains (Ethereum, Arbitrum, Optimism, Base, Polygon):
// historical price at specific block + current price.
//
// Requires API key (free tier available).
// Used to calculate PnL per trade (buy price vs current price).
//
// API docs:
//   - EVM: https://docs.moralis.com/data-api/evm/price/overview
package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"miora-ai/app/dto"
	"miora-ai/constants"
	"net/http"
)

// Moralis fetches token price data from the Moralis API.
type Moralis struct {
	apiKey string
}

// NewMoralis creates a new Moralis client with the given API key.
func NewMoralis(apiKey string) *Moralis {

	return &Moralis{apiKey: apiKey}

}

// moralisPriceResponse is the JSON response from Moralis price endpoints.
type moralisPriceResponse struct {
	UsdPrice     float64 `json:"usdPrice"`
	TokenAddress string  `json:"tokenAddress"`
}

// chainToMoralisID maps chain identifiers to Moralis EVM chain hex IDs.
func chainToMoralisID(chain string) string {

	cfg := constants.GetChainConfig(chain)
	if cfg != nil && cfg.MoralisChainID != "" {
		return cfg.MoralisChainID
	}
	return "0x1"

}

// GetTokenPrice fetches the USD price of a token.
//
// For EVM chains:
//   - If block > 0, fetches the historical price at that block.
//   - If block == 0, fetches the current price.
//   - Endpoint: GET https://deep-index.moralis.io/api/v2.2/erc20/{address}/price
func (m *Moralis) GetTokenPrice(chain, tokenAddress string, block uint64) (*dto.TokenPriceData, error) {

	return m.getEVMTokenPrice(chain, tokenAddress, block)

}

// getEVMTokenPrice fetches token price from the Moralis EVM API.
func (m *Moralis) getEVMTokenPrice(chain, tokenAddress string, block uint64) (*dto.TokenPriceData, error) {

	chainID := chainToMoralisID(chain)
	url := fmt.Sprintf(
		"https://deep-index.moralis.io/api/v2.2/erc20/%s/price?chain=%s",
		tokenAddress, chainID,
	)

	if block > 0 {
		url += fmt.Sprintf("&to_block=%d", block)
	}

	result, err := m.doRequest(url)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPriceData{
		TokenAddress: tokenAddress,
		UsdPrice:     result.UsdPrice,
		BlockNumber:  block,
	}, nil

}

// doRequest executes an authenticated GET request to Moralis and parses the response.
func (m *Moralis) doRequest(url string) (*moralisPriceResponse, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-Key", m.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moralis request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("moralis error %d: %s", resp.StatusCode, string(body))
	}

	var result moralisPriceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil

}
