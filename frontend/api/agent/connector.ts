import { request, requestVoid } from "../client";
import {
  agentConfigSchema,
  agentConfigListSchema,
  agentTradeListSchema,
  createBotRequestSchema,
  updateBotRequestSchema,
  type AgentConfig,
  type AgentTrade,
  type CreateBotRequest,
  type UpdateBotRequest,
} from "./validation";

/**
 * GET /agent/bots — List all bots for the user.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function listBots(walletAddress: string): Promise<AgentConfig[]> {
  return request("/agent/bots", agentConfigListSchema, { walletAddress });
}

/**
 * GET /agent/bots/:id — Get a single bot by ID.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function getBot(
  walletAddress: string,
  botId: number,
): Promise<AgentConfig> {
  return request(`/agent/bots/${botId}`, agentConfigSchema, { walletAddress });
}

/**
 * POST /agent/bots — Create a new bot (wallet or consensus type).
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function createBot(
  walletAddress: string,
  data: CreateBotRequest,
): Promise<AgentConfig> {
  const body = createBotRequestSchema.parse(data);
  return request("/agent/bots", agentConfigSchema, {
    method: "POST",
    body,
    walletAddress,
  });
}

/**
 * PUT /agent/bots/:id — Update bot configuration.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function updateBot(
  walletAddress: string,
  botId: number,
  data: UpdateBotRequest,
): Promise<AgentConfig> {
  const body = updateBotRequestSchema.parse(data);
  return request(`/agent/bots/${botId}`, agentConfigSchema, {
    method: "PUT",
    body,
    walletAddress,
  });
}

/**
 * DELETE /agent/bots/:id — Delete a bot.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function deleteBot(
  walletAddress: string,
  botId: number,
): Promise<void> {
  return requestVoid(`/agent/bots/${botId}`, {
    method: "DELETE",
    walletAddress,
  });
}

/**
 * POST /agent/bots/:id/start — Start a bot.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function startBot(
  walletAddress: string,
  botId: number,
): Promise<AgentConfig> {
  return request(`/agent/bots/${botId}/start`, agentConfigSchema, {
    method: "POST",
    walletAddress,
  });
}

/**
 * POST /agent/bots/:id/pause — Pause a bot.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function pauseBot(
  walletAddress: string,
  botId: number,
): Promise<AgentConfig> {
  return request(`/agent/bots/${botId}/pause`, agentConfigSchema, {
    method: "POST",
    walletAddress,
  });
}

/**
 * GET /agent/bots/:id/trades — Get bot's trade history.
 * Protected endpoint (requires X-Wallet-Address header).
 */
export async function getBotTrades(
  walletAddress: string,
  botId: number,
  limit?: number,
): Promise<AgentTrade[]> {
  const query = limit ? `?limit=${limit}` : "";
  return request(`/agent/bots/${botId}/trades${query}`, agentTradeListSchema, {
    walletAddress,
  });
}
