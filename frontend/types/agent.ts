export type AgentConfig = {
  id: number;
  user_id: number;
  budget: number;
  max_per_trade: number;
  risk_tolerance: "low" | "medium" | "high";
  min_score: number;
  conditions: string[];
  status: "active" | "paused" | "stopped";
  agent_wallet_address: string;
  total_spent: number;
  total_trades: number;
  created_at: string;
  updated_at: string;
};

export type AgentTrade = {
  id: number;
  agent_config_id: number;
  source_wallet: string;
  source_score: number;
  token_address: string;
  token_symbol: string;
  direction: "buy" | "sell";
  amount_usd: number;
  tx_hash: string;
  status: "executed" | "failed" | "skipped";
  reason: string;
  risk_assessment: string;
  created_at: string;
};
