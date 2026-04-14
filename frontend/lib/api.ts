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

// --- Swap ---

export function getSwapQuote(chain: string, inputMint: string, outputMint: string, amount: string, slippage?: number) {
  return request<import("@/types/swap").SwapQuote>("/swap/quote", {
    method: "POST",
    body: { chain, input_mint: inputMint, output_mint: outputMint, amount, slippage },
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
