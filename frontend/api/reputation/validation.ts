import { z } from "zod";

// --- Response schema (matches backend/app/dto/responses/reputation.go) ---

export const reputationSchema = z.object({
  address: z.string(),
  chain: z.string(),
  score: z.number(),
  recommendation: z.string(),
  total_transactions: z.number(),
  attestation_uid: z.string(),
  attestation_tx_hash: z.string().optional(),
  attester: z.string(),
  timestamp: z.number(),
  explorer_url: z.string(),
});

// --- Inferred types ---

export type Reputation = z.infer<typeof reputationSchema>;
