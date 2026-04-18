import type { ApiResponse } from "@/types/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

type FetchOptions = {
  method?: string;
  body?: unknown;
  token?: string;
};

/**
 * Fetch wrapper that handles the ApiResponse envelope from the backend.
 * Throws on non-success responses with the backend error message.
 */
async function request<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
  const { method = "GET", body, token } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
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

// --- Watchlist (requires auth token) ---

export function getWatchlist(token: string) {
  return request<import("@/types/watchlist").WatchlistItem[]>("/watchlist", { token });
}

export function followWallet(token: string, data: { wallet_address: string; chain: string; recommendation: string; conditions: string[]; email_notify: boolean }) {
  return request<void>("/watchlist/follow", { method: "POST", body: data, token });
}

export function unfollowWallet(token: string, address: string) {
  return request<void>(`/watchlist/${address}`, { method: "DELETE", token });
}

export function updateWatchlist(token: string, address: string, data: { conditions?: string[]; email_notify?: boolean }) {
  return request<void>(`/watchlist/${address}`, { method: "PUT", body: data, token });
}

// --- Auth ---

export function getMe(token: string) {
  return request<{ id: number; firebase_uid: string; email: string; name: string }>("/auth/me", { token });
}

// --- Reputation ---

export function getReputation(address: string) {
  return request<import("@/types/reputation").Reputation>(`/reputation/${address}`);
}

// --- Agent (requires wallet auth) ---

export function listBots(token: string) {
  return request<import("@/types/agent").AgentConfig[]>("/agent/bots", { token });
}

export function getBot(token: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}`, { token });
}

export function createBot(token: string, data: { target_wallet_address: string; target_wallet_chain: string; target_wallet_score: number; recommendation: string; budget: number; max_per_trade: number; conditions: string[] }) {
  return request<import("@/types/agent").AgentConfig>("/agent/bots", { method: "POST", body: data, token });
}

export function updateBot(token: string, botId: number, data: { budget?: number; max_per_trade?: number; conditions?: string[]; consensus_enabled?: boolean; consensus_threshold?: number; consensus_window_min?: number }) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}`, { method: "PUT", body: data, token });
}

export function deleteBot(token: string, botId: number) {
  return request<void>(`/agent/bots/${botId}`, { method: "DELETE", token });
}

export function startBot(token: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}/start`, { method: "POST", token });
}

export function pauseBot(token: string, botId: number) {
  return request<import("@/types/agent").AgentConfig>(`/agent/bots/${botId}/pause`, { method: "POST", token });
}

export function getBotTrades(token: string, botId: number, limit?: number) {
  const query = limit ? `?limit=${limit}` : "";
  return request<import("@/types/agent").AgentTrade[]>(`/agent/bots/${botId}/trades${query}`, { token });
}
