export type AgentConfig = {
  id: number;
  user_id: number;
  bot_type: "wallet" | "consensus";
  target_wallet_address?: string;
  target_wallet_chain?: string;
  target_wallet_score?: number;
  recommendation?: string;
  budget: number;
  max_per_trade: number;
  conditions: string[];
  status: "active" | "paused" | "stopped";
  agent_wallet_address: string;
  total_spent: number;
  total_trades: number;
  consensus_threshold?: number;
  consensus_window_min?: number;
  min_score?: number;
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
