package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"miora-ai/app/interfaces"
)

type AlchemyEVM struct {
	apiKey string
}

func NewAlchemyEVM(apiKey string) *AlchemyEVM {
	return &AlchemyEVM{apiKey: apiKey}
}

type alchemyEVMRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type alchemyTransfer struct {
	Hash        string  `json:"hash"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Value       float64 `json:"value"`
	Asset       string  `json:"asset"`
	BlockNum    string  `json:"blockNum"`
	RawContract struct {
		Value string `json:"value"`
	} `json:"rawContract"`
}

type alchemyEVMResponse struct {
	Result struct {
		Transfers []alchemyTransfer `json:"transfers"`
	} `json:"result"`
}

func (a *AlchemyEVM) GetTransfers(address string) ([]interfaces.TransferData, error) {
	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", a.apiKey)

	payload := alchemyEVMRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  "alchemy_getAssetTransfers",
		Params: []interface{}{
			map[string]interface{}{
				"fromAddress": address,
				"category":    []string{"external", "erc20", "erc721"},
				"order":       "desc",
				"maxCount":    "0x64",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
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
			Hash:        t.Hash,
			From:        t.From,
			To:          t.To,
			Value:       fmt.Sprintf("%f", t.Value),
			TokenSymbol: t.Asset,
			BlockNumber: hexToUint64(t.BlockNum),
		})
	}

	return transfers, nil
}
