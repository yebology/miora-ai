// Package clients provides blockchain API integrations via Alchemy.
// Each client implements the interfaces.BlockchainClient interface.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/interfaces"
	"miora-ai/constants"
)

// AlchemyEVM fetches transaction data from EVM-compatible chains (Ethereum, Arbitrum, Optimism, Base, Polygon)
// using Alchemy's alchemy_getAssetTransfers JSON-RPC method.
//
// API docs: https://docs.alchemy.com/reference/alchemy-getassettransfers
type AlchemyEVM struct {
	apiKey string
}

// NewAlchemyEVM creates a new EVM client with the given Alchemy API key.
func NewAlchemyEVM(apiKey string) *AlchemyEVM {

	return &AlchemyEVM{apiKey: apiKey}

}

// alchemyEVMRequest is the JSON-RPC request payload for Alchemy EVM endpoints.
type alchemyEVMRequest struct {
	ID      int           `json:"id"`      // JSON-RPC request identifier
	JSONRPC string        `json:"jsonrpc"` // JSON-RPC version, always "2.0"
	Method  string        `json:"method"`  // RPC method name
	Params  []interface{} `json:"params"`  // Method-specific parameters
}

// alchemyTransfer represents a single transfer from the Alchemy response.
type alchemyTransfer struct {
	Hash        string  `json:"hash"`     // Transaction hash
	From        string  `json:"from"`     // Sender address
	To          string  `json:"to"`       // Receiver address
	Value       float64 `json:"value"`    // Transfer amount (human-readable, already converted from wei)
	Asset       string  `json:"asset"`    // Token symbol (e.g. "ETH", "USDC")
	BlockNum    string  `json:"blockNum"` // Block number in hex (e.g. "0x1a4")
	RawContract struct {
		Value   string `json:"value"`   // Raw transfer value in hex (wei for ETH, smallest unit for tokens)
		Address string `json:"address"` // Token contract address
	} `json:"rawContract"` // Raw contract data, useful for precise value calculations and token identification
}

// alchemyEVMResponse is the JSON-RPC response from alchemy_getAssetTransfers.
type alchemyEVMResponse struct {
	Result struct {
		Transfers []alchemyTransfer `json:"transfers"` // List of asset transfers
	} `json:"result"`
}

// GetTransfers fetches the last 100 outgoing AND incoming transfers for the given address.
// chain parameter determines which Alchemy RPC endpoint to use (ethereum, arbitrum, optimism, base, polygon).
func (a *AlchemyEVM) GetTransfers(address string, chain ...string) ([]interfaces.TransferData, error) {

	chainKey := "ethereum"
	if len(chain) > 0 && chain[0] != "" {
		chainKey = chain[0]
	}

	cfg := constants.GetChainConfig(chainKey)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported chain: %s", chainKey)
	}

	baseURL := cfg.AlchemyURL + a.apiKey

	// Fetch outgoing transfers (wallet sent tokens = sell)
	outgoing, err := a.fetchTransfers(baseURL, address, "out")
	if err != nil {
		return nil, fmt.Errorf("fetch outgoing: %w", err)
	}

	// Fetch incoming transfers (wallet received tokens = buy)
	incoming, err := a.fetchTransfers(baseURL, address, "in")
	if err != nil {
		return nil, fmt.Errorf("fetch incoming: %w", err)
	}

	// Combine both directions
	transfers := make([]interfaces.TransferData, 0, len(outgoing)+len(incoming))
	transfers = append(transfers, outgoing...)
	transfers = append(transfers, incoming...)

	return transfers, nil

}

// fetchTransfers makes a single alchemy_getAssetTransfers call for one direction.
// direction: "in" (toAddress) or "out" (fromAddress).
func (a *AlchemyEVM) fetchTransfers(baseURL, address, direction string) ([]interfaces.TransferData, error) {

	params := map[string]interface{}{
		"category": []string{"external", "erc20", "erc721"},
		"order":    "desc",
		"maxCount": "0x64",
	}

	if direction == "out" {
		params["fromAddress"] = address
	} else {
		params["toAddress"] = address
	}

	payload := alchemyEVMRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  "alchemy_getAssetTransfers",
		Params:  []interface{}{params},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("alchemy evm request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result alchemyEVMResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	transfers := make([]interfaces.TransferData, 0, len(result.Result.Transfers))
	for _, t := range result.Result.Transfers {
		transfers = append(transfers, interfaces.TransferData{
			Hash:            t.Hash,
			From:            t.From,
			To:              t.To,
			Value:           fmt.Sprintf("%f", t.Value),
			TokenSymbol:     t.Asset,
			ContractAddress: t.RawContract.Address,
			BlockNumber:     hexToUint64(t.BlockNum),
			Direction:       direction,
		})
	}

	return transfers, nil

}
