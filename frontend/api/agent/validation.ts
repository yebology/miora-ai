import { z } from "zod";

// --- Response schemas (matches backend/app/entities/agent_config.go + agent_trade.go) ---

export const agentConfigSchema = z.object({
  id: z.number(),
  user_id: z.number(),
  bot_type: z.enum(["wallet", "consensus"]),
  target_wallet_address: z.string().optional(),
  target_wallet_chain: z.string().optional(),
  target_wallet_score: z.number().optional(),
  recommendation: z.string().optional(),
  budget: z.number(),
  max_per_trade: z.number(),
  conditions: z.union([z.array(z.string()), z.null()]).transform((v) => v ?? []),
  status: z.enum(["active", "paused", "stopped"]),
  agent_wallet_address: z.string(),
  total_spent: z.number(),
  total_trades: z.number(),
  consensus_threshold: z.number().optional(),
  consensus_window_min: z.number().optional(),
  min_score: z.number().optional(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const agentConfigListSchema = z.array(agentConfigSchema);

export const agentTradeSchema = z.object({
  id: z.number(),
  agent_config_id: z.number(),
  source_wallet: z.string(),
  source_score: z.number(),
  token_address: z.string(),
  token_symbol: z.string(),
  direction: z.enum(["buy", "sell"]),
  amount_usd: z.number(),
  tx_hash: z.string(),
  status: z.enum(["executed", "failed", "skipped"]),
  reason: z.string(),
  risk_assessment: z.string(),
  created_at: z.string(),
});

export const agentTradeListSchema = z.array(agentTradeSchema);

// --- Request validation schemas (matches backend/app/dto/requests/agent.go) ---

export const createBotRequestSchema = z.object({
  bot_type: z.enum(["wallet", "consensus"]),
  target_wallet_address: z.string().optional(),
  target_wallet_chain: z.string().optional(),
  target_wallet_score: z.number().optional(),
  recommendation: z.string().optional(),
  budget: z.number().positive("Budget must be greater than 0"),
  max_per_trade: z.number().positive("Max per trade must be greater than 0"),
  conditions: z.array(z.string()).default([]),
  consensus_threshold: z.number().int().min(2).max(20).optional(),
  consensus_window_min: z.number().int().min(5).max(1440).optional(),
  min_score: z.number().int().min(0).max(100).optional(),
});

export const updateBotRequestSchema = z.object({
  budget: z.number().min(0).optional(),
  max_per_trade: z.number().min(0).optional(),
  conditions: z.array(z.string()).optional(),
  consensus_threshold: z.number().int().min(2).max(20).optional(),
  consensus_window_min: z.number().int().min(5).max(1440).optional(),
  min_score: z.number().int().min(0).max(100).optional(),
});

// --- Inferred types ---

export type AgentConfig = z.infer<typeof agentConfigSchema>;
export type AgentTrade = z.infer<typeof agentTradeSchema>;
export type CreateBotRequest = z.infer<typeof createBotRequestSchema>;
export type UpdateBotRequest = z.infer<typeof updateBotRequestSchema>;
