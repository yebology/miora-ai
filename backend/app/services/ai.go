package services

import (
	"miora-ai/app/dto/prompts"
	"miora-ai/app/dto/responses"
	"miora-ai/app/interfaces"
)

// AIService generates natural language insights from wallet scoring data.
type AIService struct {
	ai interfaces.IAI
}

// NewAIService creates a new AIService.
func NewAIService(ai interfaces.IAI) *AIService {

	return &AIService{ai: ai}

}

// GenerateInsight takes wallet analysis data and returns a natural language explanation.
// Tone defaults to "simple" if empty.
func (s *AIService) GenerateInsight(analysis *responses.WalletAnalysis, tone string) (string, error) {

	if tone == "" {
		tone = "simple"
	}
	prompt := prompts.BuildWalletInsight(analysis, tone)
	return s.ai.Generate(prompt)

}

// GenerateCustomInsight generates an insight using a user-provided custom prompt.
func (s *AIService) GenerateCustomInsight(analysis *responses.WalletAnalysis, customPrompt string) (string, error) {

	prompt := prompts.BuildWalletInsightCustom(analysis, customPrompt)
	return s.ai.Generate(prompt)

}

// GenerateTradeAssessment generates a short AI risk assessment for a new trade notification.
func (s *AIService) GenerateTradeAssessment(walletAddress, chain, tokenSymbol, direction string, liquidity, marketCap, priceChange24h, pairAgeHours float64) (string, error) {

	prompt := prompts.BuildTradeAssessment(walletAddress, chain, tokenSymbol, direction, liquidity, marketCap, priceChange24h, pairAgeHours)
	return s.ai.Generate(prompt)

}
