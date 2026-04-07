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
func (s *AIService) GenerateInsight(analysis *responses.WalletAnalysis) (string, error) {

	prompt := prompts.BuildWalletInsight(analysis)
	return s.ai.Generate(prompt)

}
