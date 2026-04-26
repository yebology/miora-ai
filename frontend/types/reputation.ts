export type Reputation = {
  address: string;
  chain: string;
  score: number;
  recommendation: string;
  total_transactions: number;
  attestation_uid: string;
  attestation_tx_hash?: string;
  attester?: string;
  timestamp?: number;
  explorer_url: string;
};
