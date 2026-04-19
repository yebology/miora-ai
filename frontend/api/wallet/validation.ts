import { z } from "zod";

// --- Response schemas (matches backend/app/dto/responses/wallet.go) ---

export const tradedTokenSchema = z.object({
  contract_address: z.string(),
  symbol: z.string(),
  chain: z.string(),
  pnl_percent: z.number(),
  buy_price: z.number(),
  exit_price: z.number(),
  buy_time: z.string(),
  exit_time: z.string().nullable().optional().transform((v) => v ?? undefined),
  status: z.enum(["realized", "unrealized"]),
});

export const conditionSchema = z.object({
  id: z.string(),
  label: z.string(),
  description: z.string(),
  type: z.string(),
  field: z.string(),
  operator: z.string(),
  value: z.coerce.number(),
});

export const walletAnalysisSchema = z.object({
  address: z.string(),
  chain: z.string(),
  total_transactions: z.number(),
  profit_consistency: z.number(),
  win_rate: z.number(),
  risk_exposure: z.number(),
  entry_timing: z.number(),
  token_quality: z.number(),
  trade_discipline: z.number(),
  final_score: z.number(),
  recommendation: z.enum(["full_follow", "conditional_follow", "avoid"]),
  ai_insight: z.string().optional(),
  ai_insight_tone: z.string().optional(),
  ai_insight_prompt: z.string().optional(),
  traded_tokens: z.array(tradedTokenSchema).optional(),
  conditions: z.array(conditionSchema).optional(),
});

export const regenerateInsightSchema = z.object({
  ai_insight: z.string(),
  tone: z.string(),
});

// --- Request validation schemas ---

export const analyzeWalletRequestSchema = z.object({
  address: z.string().min(1, "Address is required"),
  chain: z.string().min(1, "Chain is required"),
  limit: z.number().int().positive().optional(),
});

export const regenerateInsightRequestSchema = z.object({
  address: z.string().min(1, "Address is required"),
  chain: z.string().min(1, "Chain is required"),
  tone: z.enum(["simple", "eli5", "custom"]),
  custom_prompt: z.string().optional(),
});

// --- Inferred types ---

export type WalletAnalysis = z.infer<typeof walletAnalysisSchema>;
export type TradedToken = z.infer<typeof tradedTokenSchema>;
export type Condition = z.infer<typeof conditionSchema>;
export type RegenerateInsightResponse = z.infer<typeof regenerateInsightSchema>;
export type AnalyzeWalletRequest = z.infer<typeof analyzeWalletRequestSchema>;
export type RegenerateInsightRequest = z.infer<typeof regenerateInsightRequestSchema>;
