import { request } from "../client";
import {
  walletAnalysisSchema,
  regenerateInsightSchema,
  analyzeWalletRequestSchema,
  regenerateInsightRequestSchema,
  type WalletAnalysis,
  type RegenerateInsightResponse,
} from "./validation";

/**
 * POST /wallets/analyze — Analyze a wallet on Base.
 * Public endpoint (no auth required).
 */
export async function analyzeWallet(
  address: string,
  chain: string,
  limit?: number,
): Promise<WalletAnalysis> {
  const body = analyzeWalletRequestSchema.parse({ address, chain, limit });
  return request("/wallets/analyze", walletAnalysisSchema, {
    method: "POST",
    body,
  });
}

/**
 * GET /wallets/:address — Get stored analysis for a wallet.
 * Public endpoint (no auth required).
 */
export async function getWallet(address: string): Promise<WalletAnalysis> {
  return request(`/wallets/${address}`, walletAnalysisSchema);
}

/**
 * POST /wallets/regenerate-insight — Regenerate AI insight with a different tone.
 * Public endpoint (no auth required).
 */
export async function regenerateInsight(
  address: string,
  chain: string,
  tone: "simple" | "eli5" | "custom",
  customPrompt?: string,
): Promise<RegenerateInsightResponse> {
  const body = regenerateInsightRequestSchema.parse({
    address,
    chain,
    tone,
    custom_prompt: customPrompt,
  });
  return request("/wallets/regenerate-insight", regenerateInsightSchema, {
    method: "POST",
    body,
  });
}
