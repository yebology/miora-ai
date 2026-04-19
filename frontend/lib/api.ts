import type { ApiResponse } from "@/types/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

type FetchOptions = {
  method?: string;
  body?: unknown;
  walletAddress?: string;
};

/**
 * Fetch wrapper that handles the ApiResponse envelope from the backend.
 * Throws on non-success responses with the backend error message.
 *
 * For protected routes, pass walletAddress — sent as X-Wallet-Address header.
 */
async function request<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
  const { method = "GET", body, walletAddress } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (walletAddress) {
    headers["X-Wallet-Address"] = walletAddress;
  }

  const res = await fetch(`${API_URL}${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const json: ApiResponse<T> = await res.json();

  if (json.status === "error" || !res.ok) {
    throw new Error(json.message || "Something went wrong.");
  }

  return json.data as T;
}

// --- Wallet ---

export function analyzeWallet(address: string, chain: string, limit: number) {
  return request<import("@/types/wallet").WalletAnalysis>("/wallets/analyze", {
    method: "POST",
    body: { address, chain, limit },
  });
}

export function getWallet(address: string) {
  return request<import("@/types/wallet").WalletAnalysis>(`/wallets/${address}`);
}

export function regenerateInsight(address: string, chain: string, tone: string, customPrompt?: string) {
  return request<{ ai_insight: string; tone: string }>("/wallets/regenerate-insight", {
    method: "POST",
    body: { address, chain, tone, custom_prompt: customPrompt },
  });
}

// --- Watchlist (requires wallet auth — X-Wallet-Address header) ---

export function getWatchlist(walletAddress: string) {
  return request<import("@/types/watchlist").WatchlistItem[]>("/watchlist", { walletAddress });
}

export function followWallet(walletAddress: string, data: { wallet_address: string; chain: string; recommendation: string; conditions: string[]; email_notify: boolean }) {
  return request<void>("/watchlist/follow", { method: "POST", body: data, walletAddress });
}

export function unfollowWallet(walletAddress: string, address: string) {
  return request<void>(`/watchlist/${address}`, { method: "DELETE", walletAddress });
}

export function updateWatchlist(walletAddress: string, address: string, data: { conditions?: string[]; email_notify?: boolean }) {
  return request<void>(`/watchlist/${address}`, { method: "PUT", body: data, walletAddress });
}

// --- Auth (requires wallet auth) ---

export function getMe(walletAddress: string) {
  return request<{ id: number; wallet_address: string }>("/auth/me", { walletAddress });
}

// --- Reputation ---

export function getReputation(address: string) {
  return request<import("@/types/reputation").Reputation>(`/reputation/${address}`);
}

// --- Agent (requires wallet auth — X-Wallet-Address header) ---

export function listBots(walletAddress: string) {
  return request<import("@/types/agent").AgentConfig[]>("/agent/bots", { walletAddress });
}

export function getBot(walletAddress: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}`, { walletAddress });
}

export function createBot(walletAddress: string, data: { target_wallet_address: string; target_wallet_chain: string; target_wallet_score: number; recommendation: string; budget: number; max_per_trade: number; conditions: string[] }) {
  return request<import("@/types/agent").AgentConfig>("/agent/bots", { method: "POST", body: data, walletAddress });
}

export function updateBot(walletAddress: string, botId: number, data: { budget?: number; max_per_trade?: number; conditions?: string[]; consensus_enabled?: boolean; consensus_threshold?: number; consensus_window_min?: number }) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}`, { method: "PUT", body: data, walletAddress });
}

export function deleteBot(walletAddress: string, botId: number) {
  return request<void>(`/agent/bots/${botId}`, { method: "DELETE", walletAddress });
}

export function startBot(walletAddress: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}/start`, { method: "POST", walletAddress });
}

export function pauseBot(walletAddress: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}/pause`, { method: "POST", walletAddress });
}

export function getBotTrades(walletAddress: string, botId: number, limit?: number) {
  const query = limit ? `?limit=${limit}` : "";
  return request<import("@/types/agent").AgentTrade[]>(`/agent/bots/${botId}/trades${query}`, { walletAddress });
}
