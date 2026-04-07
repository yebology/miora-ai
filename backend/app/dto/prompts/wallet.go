package prompts

import (
	"fmt"

	"miora-ai/app/dto/responses"
)

// BuildWalletInsight constructs the LLM prompt for wallet analysis insight.
// Tone controls the language style: "simple" (default), "technical", or "eli5".
func BuildWalletInsight(a *responses.WalletAnalysis, tone string) string {

	toneInstruction := getToneInstruction(tone)

	return fmt.Sprintf(`You are Miora AI, a blockchain wallet analyst.

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

%s

- Classify the trading style in simple terms (e.g. "quick flipper", "patient investor", "risky gambler")
- End with a clear recommendation: should someone follow this wallet or not, and why
- Do not use markdown formatting, bullet points, or headers — just plain text`,
		a.Address, a.Chain, a.TotalTransactions,
		a.WinRate, a.ProfitConsistency, a.EntryTiming,
		a.TokenQuality, a.TradeDiscipline, a.RiskExposure,
		a.FinalScore, a.Recommendation,
		toneInstruction,
	)

}

// getToneInstruction returns prompt instructions based on the selected tone.
func getToneInstruction(tone string) string {

	switch tone {
	case "eli5":
		return `Tone: Explain Like I'm 5
- Use very simple words and fun analogies (like comparing trading to a lemonade stand or collecting cards)
- Avoid ALL crypto and financial terms
- Make it feel like a story a kid would understand
- Keep it light and friendly`
	default:
		return `Tone: Simple (Beginner-Friendly)
- Use simple, everyday language — no crypto jargon
- Explain what this wallet does well and what it doesn't
- Write as if explaining to someone who just started learning about crypto`
	}

}

// BuildWalletInsightCustom constructs the LLM prompt with a user-provided custom instruction.
func BuildWalletInsightCustom(a *responses.WalletAnalysis, customPrompt string) string {

	return fmt.Sprintf(`You are Miora AI, a blockchain wallet analyst.

Analyze this wallet based on the user's specific request.

Wallet: %s
Chain: %s
Total Transactions: %d

Scoring (0-100, higher is better):
- Win Rate: %.2f
- Profit Consistency: %.2f
- Entry Timing: %.2f
- Token Quality: %.2f
- Trade Discipline: %.2f
- Risk Exposure: %.2f

Final Score: %.2f out of 100
Recommendation: %s

User's request: "%s"

Instructions:
- Answer the user's specific request using the wallet data above
- Keep it concise (3-5 sentences)
- Do not use markdown formatting, bullet points, or headers — just plain text
- Stay factual — only reference data provided above`,
		a.Address, a.Chain, a.TotalTransactions,
		a.WinRate, a.ProfitConsistency, a.EntryTiming,
		a.TokenQuality, a.TradeDiscipline, a.RiskExposure,
		a.FinalScore, a.Recommendation,
		customPrompt,
	)

}
