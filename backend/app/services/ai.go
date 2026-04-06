package services

import (
	"fmt"

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

	prompt := buildPrompt(analysis)
	return s.ai.Generate(prompt)

}

// buildPrompt constructs the LLM prompt from wallet scoring data.
func buildPrompt(a *responses.WalletAnalysis) string {

	return fmt.Sprintf(`You are Miora AI, a friendly blockchain wallet analyst that explains things simply.
Your audience is beginners who may not understand crypto jargon.

Analyze this wallet and write a short, clear explanation (3-4 sentences).

Wallet: %s
Chain: %s
Total Transactions: %d

Scoring (0-100, higher is better):
- Win Rate: %.2f (how often trades are profitable)
- Profit Consistency: %.2f (how stable the profits are)
- Entry Timing: %.2f (how early they enter new tokens)
- Token Quality: %.2f (how reputable the tokens they trade are)
- Trade Discipline: %.2f (how focused their trading is)
- Risk Exposure: %.2f (percentage of risky/low-liquidity tokens)

Final Score: %.2f out of 100
Recommendation: %s

Instructions:
- Use simple, everyday language — no crypto jargon
- Explain what this wallet does well and what it doesn't
- Classify the trading style in simple terms (e.g. "quick flipper", "patient investor", "risky gambler")
- End with a clear recommendation: should someone follow this wallet or not, and why
- Write as if explaining to someone who just started learning about crypto
- Do not use markdown formatting, bullet points, or headers — just plain text`,
		a.Address, a.Chain, a.TotalTransactions,
		a.WinRate, a.ProfitConsistency, a.EntryTiming,
		a.TokenQuality, a.TradeDiscipline, a.RiskExposure,
		a.FinalScore, a.Recommendation,
	)

}
