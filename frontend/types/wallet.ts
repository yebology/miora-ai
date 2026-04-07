export type TradedToken = {
  contract_address: string;
  symbol: string;
  chain: string;
  pnl_percent: number;
  buy_price: number;
  exit_price: number;
  buy_time: string;
  exit_time?: string;
  status: "realized" | "unrealized";
};

export type Condition = {
  id: string;
  label: string;
  type: string;
  field: string;
  operator: string;
  value: number;
};

export type WalletAnalysis = {
  address: string;
  chain: string;
  total_transactions: number;
  profit_consistency: number;
  win_rate: number;
  risk_exposure: number;
  entry_timing: number;
  token_quality: number;
  trade_discipline: number;
  final_score: number;
  recommendation: "full_follow" | "conditional_follow" | "avoid";
  ai_insight?: string;
  traded_tokens?: TradedToken[];
  conditions?: Condition[];
};
