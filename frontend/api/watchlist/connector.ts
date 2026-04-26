import { request, requestVoid } from "../client";
import {
  watchlistListSchema,
  followWalletRequestSchema,
  updateWatchlistRequestSchema,
  type WatchlistItem,
  type FollowWalletRequest,
  type UpdateWatchlistRequest,
} from "./validation";

/**
 * GET /watchlist — List all followed wallets for the user.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function getWatchlist(
  walletAddress: string,
): Promise<WatchlistItem[]> {
  return request("/watchlist", watchlistListSchema, { walletAddress });
}

/**
 * POST /watchlist/follow — Follow a wallet with conditions.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function followWallet(
  walletAddress: string,
  data: FollowWalletRequest,
): Promise<void> {
  const body = followWalletRequestSchema.parse(data);
  return requestVoid("/watchlist/follow", {
    method: "POST",
    body,
    walletAddress,
  });
}

/**
 * DELETE /watchlist/:address — Unfollow a wallet.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function unfollowWallet(
  walletAddress: string,
  targetAddress: string,
): Promise<void> {
  return requestVoid(`/watchlist/${targetAddress}`, {
    method: "DELETE",
    walletAddress,
  });
}

/**
 * PUT /watchlist/:address — Update conditions or notification preference.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function updateWatchlist(
  walletAddress: string,
  targetAddress: string,
  data: UpdateWatchlistRequest,
): Promise<void> {
  const body = updateWatchlistRequestSchema.parse(data);
  return requestVoid(`/watchlist/${targetAddress}`, {
    method: "PUT",
    body,
    walletAddress,
  });
}
