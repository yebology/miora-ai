// Package clients provides blockchain API integrations via Alchemy.
// Each client implements the interfaces.BlockchainClient interface.
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

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
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// alchemyTransfer represents a single transfer from the Alchemy response.
type alchemyTransfer struct {
	Hash        string  `json:"hash"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Value       float64 `json:"value"`
	Asset       string  `json:"asset"`
	BlockNum    string  `json:"blockNum"`
	RawContract struct {
		Value   string `json:"value"`
		Address string `json:"address"`
	} `json:"rawContract"`
}

// alchemyEVMResponse is the JSON-RPC response from alchemy_getAssetTransfers.
type alchemyEVMResponse struct {
	Result struct {
		Transfers []alchemyTransfer `json:"transfers"`
	} `json:"result"`
}

// GetTransfers fetches the last 100 outgoing AND incoming transfers for the given address.
func (a *AlchemyEVM) GetTransfers(address string, limit int, chain ...string) ([]interfaces.TransferData, error) {

	chainKey := "ethereum"
	if len(chain) > 0 && chain[0] != "" {
		chainKey = chain[0]
	}

	cfg := constants.GetChainConfig(chainKey)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported chain: %s", chainKey)
	}

	if limit <= 0 || limit > 50 {
		limit = 25
	}

	baseURL := cfg.AlchemyURL + a.apiKey

	// Fetch limit transfers per direction, then merge and take the most recent `limit`
	outgoing, err := a.fetchTransfers(baseURL, address, "out", limit, cfg.BlockTimeSec)
	if err != nil {
		return nil, fmt.Errorf("fetch outgoing: %w", err)
	}

	incoming, err := a.fetchTransfers(baseURL, address, "in", limit, cfg.BlockTimeSec)
	if err != nil {
		return nil, fmt.Errorf("fetch incoming: %w", err)
	}

	// Combine, sort by block number descending, take top `limit`
	all := make([]interfaces.TransferData, 0, len(outgoing)+len(incoming))
	all = append(all, outgoing...)
	all = append(all, incoming...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].BlockNumber > all[j].BlockNumber
	})

	if len(all) > limit {
		all = all[:limit]
	}

	return all, nil

}

// fetchTransfers makes a single alchemy_getAssetTransfers call for one direction.
func (a *AlchemyEVM) fetchTransfers(baseURL, address, direction string, limit int, blockTimeSec float64) ([]interfaces.TransferData, error) {

	params := map[string]interface{}{
		"category": []string{"external", "erc20", "erc721"},
		"order":    "desc",
		"maxCount": fmt.Sprintf("0x%x", limit),
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
		blockNum := hexToUint64(t.BlockNum)
		transfers = append(transfers, interfaces.TransferData{
			Hash:            t.Hash,
			From:            t.From,
			To:              t.To,
			Value:           fmt.Sprintf("%f", t.Value),
			TokenSymbol:     t.Asset,
			ContractAddress: t.RawContract.Address,
			BlockNumber:     blockNum,
			Timestamp:       estimateBlockTimestamp(blockNum, blockTimeSec),
			Direction:       direction,
		})
	}

	return transfers, nil

}

// estimateBlockTimestamp estimates a unix timestamp from a block number.
// Uses chain-specific block time (seconds per block).
//
// For Ethereum mainnet (12s/block): uses The Merge as reference point.
// For L2s (Arbitrum 0.25s, Optimism/Base/Polygon 2s): uses genesis-based estimation.
//
// This is an approximation — not exact, but close enough for display.
func estimateBlockTimestamp(blockNum uint64, blockTimeSec float64) int64 {

	if blockTimeSec == 0 {
		return 0
	}

	// Ethereum mainnet reference point
	if blockTimeSec == 12 {
		const mergeBlock uint64 = 15537394
		const mergeTimestamp int64 = 1663224162 // Sep 15, 2022
		if blockNum >= mergeBlock {
			return mergeTimestamp + int64(float64(blockNum-mergeBlock)*blockTimeSec)
		}
		return mergeTimestamp - int64(float64(mergeBlock-blockNum)*13)
	}

	// L2 chains: estimate from block number × block time
	// Most L2s started relatively recently, so block 0 ≈ chain launch
	// Arbitrum genesis: Mar 2023, Optimism: Jun 2022, Base: Aug 2023, Polygon: May 2020
	return int64(float64(blockNum) * blockTimeSec)

}
