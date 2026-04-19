import { z } from "zod";

// --- Response schema (matches backend/app/entities/watchlist.go) ---

export const watchlistItemSchema = z.object({
  id: z.number(),
  user_id: z.number(),
  wallet_address: z.string(),
  chain: z.string(),
  recommendation: z.string(),
  conditions: z.union([z.array(z.string()), z.null()]).transform((v) => v ?? []),
  email_notify: z.boolean(),
  created_at: z.string(),
});

export const watchlistListSchema = z.array(watchlistItemSchema);

// --- Request validation schemas ---

export const followWalletRequestSchema = z.object({
  wallet_address: z.string().min(1, "Wallet address is required"),
  chain: z.string().min(1, "Chain is required"),
  recommendation: z.string(),
  conditions: z.array(z.string()),
  email_notify: z.boolean(),
});

export const updateWatchlistRequestSchema = z.object({
  conditions: z.array(z.string()).optional(),
  email_notify: z.boolean().optional(),
});

// --- Inferred types ---

export type WatchlistItem = z.infer<typeof watchlistItemSchema>;
export type FollowWalletRequest = z.infer<typeof followWalletRequestSchema>;
export type UpdateWatchlistRequest = z.infer<typeof updateWatchlistRequestSchema>;
