import { request } from "../client";
import { reputationSchema, type Reputation } from "./validation";

/**
 * GET /reputation/:address — Get on-chain reputation attestation for a wallet.
 * Public endpoint (no auth required).
 */
export async function getReputation(address: string): Promise<Reputation> {
  return request(`/reputation/${address}`, reputationSchema);
}
