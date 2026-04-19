import { z } from "zod";

/**
 * Core API client — shared by all module connectors.
 *
 * Handles:
 * - Request construction (method, body, X-Wallet-Address header)
 * - Response envelope validation ({ status, message, data })
 * - Data payload validation via Zod schema
 */

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

const apiResponseSchema = z.object({
  status: z.enum(["success", "error"]),
  message: z.string(),
  data: z.unknown().optional(),
});

type RequestOptions = {
  method?: string;
  body?: unknown;
  walletAddress?: string;
};

/**
 * Sends a request and validates the response data with the given Zod schema.
 */
export async function request<T>(
  endpoint: string,
  schema: z.ZodType<T>,
  options: RequestOptions = {},
): Promise<T> {
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

  const raw = await res.json();
  const envelope = apiResponseSchema.parse(raw);

  if (envelope.status === "error" || !res.ok) {
    throw new Error(envelope.message || "Something went wrong.");
  }

  return schema.parse(envelope.data);
}

/**
 * Sends a request for endpoints that return no data payload (DELETE, follow, etc.).
 */
export async function requestVoid(
  endpoint: string,
  options: RequestOptions = {},
): Promise<void> {
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

  const raw = await res.json();
  const envelope = apiResponseSchema.parse(raw);

  if (envelope.status === "error" || !res.ok) {
    throw new Error(envelope.message || "Something went wrong.");
  }
}
