package services

import (
	"miora-ai/app/dto/responses"
	"miora-ai/app/interfaces"
	"miora-ai/constants"
	"miora-ai/pkg"
)

// SwapService handles swap quote requests by routing to the correct DEX aggregator.
type SwapService struct {
	jupiter interfaces.ISwapClient
	oneInch interfaces.ISwapClient
}

// NewSwapService creates a new SwapService.
func NewSwapService(jupiter, oneInch interfaces.ISwapClient) *SwapService {

	return &SwapService{
		jupiter: jupiter,
		oneInch: oneInch,
	}

}

// GetQuote fetches a swap quote from the appropriate DEX aggregator.
// Solana → Jupiter, EVM chains → 1inch.
func (s *SwapService) GetQuote(chain, inputMint, outputMint, amount string, slippage int) (*responses.SwapQuote, *pkg.AppError) {

	cfg := constants.GetChainConfig(chain)
	if cfg == nil {
		return nil, pkg.ErrBadReq(constants.UnsupportedChain)
	}

	if constants.IsSolana(chain) {
		quote, err := s.jupiter.GetQuote(inputMint, outputMint, amount, slippage)
		if err != nil {
			return nil, pkg.ErrUnexpected(502, err.Error())
		}
		return quote, nil
	}

	if constants.IsEVM(chain) {
		quote, err := s.oneInch.GetQuote(inputMint, outputMint, amount, slippage, chain)
		if err != nil {
			return nil, pkg.ErrUnexpected(502, err.Error())
		}
		quote.Chain = chain
		return quote, nil
	}

	return nil, pkg.ErrBadReq(constants.UnsupportedChain)

}
