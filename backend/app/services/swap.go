package services

import (
	"miora-ai/app/dto/responses"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"
)

// SwapService handles swap quote requests by routing to 1inch.
type SwapService struct {
	oneInch interfaces.ISwapClient
}

// NewSwapService creates a new SwapService.
func NewSwapService(oneInch interfaces.ISwapClient) *SwapService {

	return &SwapService{
		oneInch: oneInch,
	}

}

// GetQuote fetches a swap quote from 1inch for EVM chains.
func (s *SwapService) GetQuote(chain, inputMint, outputMint, amount string, slippage int) (*responses.SwapQuote, *pkg.AppError) {

	cfg := constants.GetChainConfig(chain)
	if cfg == nil {
		return nil, pkg.ErrBadReq(constants.UnsupportedChain)
	}

	if !constants.IsEVM(chain) {
		return nil, pkg.ErrBadReq(constants.UnsupportedChain)
	}

	quote, err := s.oneInch.GetQuote(inputMint, outputMint, amount, slippage, chain)
	if err != nil {
		return nil, pkg.ErrUnexpected(502, err.Error())
	}
	quote.Chain = chain
	return quote, nil

}
