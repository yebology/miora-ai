export type WatchlistItem = {
  id: number;
  wallet_address: string;
  chain: string;
  recommendation: string;
  conditions: string[]; // selected condition IDs
  email_notify: boolean;
  created_at: string;
};

export type Notification = {
  id: number;
  wallet_address: string;
  chain: string;
  token_address: string;
  token_symbol: string;
  direction: "in" | "out"; // in = buy, out = sell
  value: string;
  liquidity: number;
  market_cap: number;
  read: boolean;
  created_at: string;
};
