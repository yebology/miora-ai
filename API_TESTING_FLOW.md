# Miora AI — API Testing Flow

## Prerequisites

Before testing, make sure everything is ready:

1. Fill in all API keys in `backend/.env`
   - `ALCHEMY_API_KEY` — from [Alchemy Dashboard](https://dashboard.alchemy.com/)
   - `MORALIS_API_KEY` — from [Moralis Admin](https://admin.moralis.io/)
   - `BIRDEYE_API_KEY` — from [Birdeye](https://birdeye.so/)
   - `GEMINI_API_KEY` — from [Google AI Studio](https://aistudio.google.com/)
   - `ONEINCH_API_KEY` — from [1inch Dev Portal](https://portal.1inch.dev/)
   - `RESEND_API_KEY` — from [Resend Dashboard](https://resend.com/) (optional, for email notifications)
2. Place `firebase-credentials.json` in the `backend/` folder — download from Firebase Console > Project Settings > Service Accounts
3. Start database: `cd backend && docker compose up -d`
4. Start server: `cd backend && go run main.go`
5. Install [Bruno](https://www.usebruno.com/) and import collection from `backend/api-docs/`

Related files:
- Config loader: `backend/config/config.go`
- Docker Compose: `backend/docker-compose.yml`
- Entry point: `backend/main.go`
- Route setup: `backend/router/routes.go`
- DI container: `backend/router/container.go`

---

## Step 1 — Health Check

**Endpoint:** `GET /api/health`
**Auth:** Not required

```
GET http://localhost:8080/api/health
```

**What to verify:**
- Response is `{"status": "ok"}`
- If it fails, the server isn't running or the port is wrong

Related files: `backend/router/routes.go` (health route)

---

## Step 2 — Analyze Wallet (EVM)

**Endpoint:** `POST /api/wallets/analyze`
**Auth:** Not required

```json
{
  "address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
  "chain": "ethereum",
  "limit": 10
}
```

**Valid EVM limits:** 10, 25, 50, 100 (default: 10)

**What to verify in response:**
- `status` = `"success"`
- `data.address` = the address you sent
- `data.chain` = `"ethereum"`
- `data.total_transactions` > 0
- Scoring metrics (all 0-100):
  - `data.win_rate` — percentage of profitable trades
  - `data.profit_consistency` — profit stability across trades
  - `data.entry_timing` — how early the wallet enters new tokens
  - `data.token_quality` — average market cap of traded tokens
  - `data.trade_discipline` — trading focus (unique tokens vs total tx)
  - `data.risk_exposure` — percentage of low-liquidity tokens (informational, not in score formula)
- `data.final_score` — overall score (0-100)
- `data.recommendation` — one of: `"full_follow"`, `"conditional_follow"`, `"avoid"`
- `data.traded_tokens` — array of traded tokens, verify:
  - `contract_address` is not empty
  - `pnl_percent` exists (can be negative)
  - `status` = `"realized"` or `"unrealized"`
  - `buy_price` and `exit_price` > 0
- `data.ai_insight` — text insight from Gemini (empty if API key not set)
- If `recommendation` = `"conditional_follow"`, check `data.conditions`:
  - Each condition has `id`, `label`, `description`, `field`, `operator`, `value`
  - Example IDs: `"min_liquidity"`, `"min_pair_age"`, `"min_mcap"`, `"min_volume"`

Related files:
- Request DTO: `backend/app/dto/requests/wallet.go`
- Response DTO: `backend/app/dto/responses/wallet.go`
- Handler: `backend/app/handlers/wallet.go`
- Service: `backend/app/services/wallet.go`
- Scoring: `backend/app/services/scoring.go`
- AI insight: `backend/app/services/ai.go`
- Limit config: `backend/constants/limits.go`
- EVM client: `backend/app/clients/evm.go`

---

## Step 3 — Analyze Wallet (Solana)

**Endpoint:** `POST /api/wallets/analyze`

```json
{
  "address": "JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4",
  "chain": "solana",
  "limit": 20
}
```

**Valid Solana limits:** 20, 50, 100, 200 (default: 20)

**What to verify:** Same as Step 2, but also check:
- `data.chain` = `"solana"`
- Solana uses Birdeye for historical price (not Moralis)
- Token addresses are base58 format (not 0x)

Related files:
- Solana client: `backend/app/clients/solana.go`
- Birdeye client: `backend/app/clients/birdeye.go`

---

## Step 4 — Get Stored Wallet

**Endpoint:** `GET /api/wallets/:address`
**Auth:** Not required

```
GET http://localhost:8080/api/wallets/0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

**What to verify:**
- Response contains the same data as the previous analyze result
- `traded_tokens` and `conditions` are NOT included (only scoring + recommendation)
- If the address was never analyzed, returns 404

Related files: `backend/app/handlers/wallet.go` (GetWallet method)

---

## Step 5 — Regenerate AI Insight

**Endpoint:** `POST /api/wallets/regenerate-insight`
**Auth:** Not required

```json
{
  "address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
  "chain": "ethereum",
  "tone": "simple"
}
```

**Tone options:** `"simple"`, `"eli5"`, `"custom"`

For custom tone:
```json
{
  "address": "0x...",
  "chain": "ethereum",
  "tone": "custom",
  "custom_prompt": "Explain in simple terms for a beginner"
}
```

**What to verify:**
- `data.ai_insight` — new text different from the previous insight
- `data.tone` — matches what you sent
- Wallet must have been analyzed before (otherwise → 404)

Related files:
- Request DTO: `backend/app/dto/requests/insight.go`
- AI service: `backend/app/services/ai.go`
- Gemini client: `backend/app/clients/gemini.go`

---

## Step 6 — Swap Quote (Solana via Jupiter)

**Endpoint:** `POST /api/swap/quote`
**Auth:** Not required

```json
{
  "chain": "solana",
  "input_mint": "So11111111111111111111111111111111111111112",
  "output_mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
  "amount": "100000000",
  "slippage": 50
}
```

**What to verify:**
- `data.chain` = `"svm"` or `"solana"`
- `data.input_amount` = the amount you sent
- `data.output_amount` — amount of tokens to receive
- `data.price_impact` — price impact percentage (lower is better)
- `data.route` — DEX route used (e.g. "Raydium → Orca")

Related files:
- Request DTO: `backend/app/dto/requests/swap.go`
- Response DTO: `backend/app/dto/responses/swap.go`
- Handler: `backend/app/handlers/swap.go`
- Service: `backend/app/services/swap.go`
- Jupiter client: `backend/app/clients/jupiter.go`

---

## Step 7 — Swap Quote (EVM via 1inch)

**Endpoint:** `POST /api/swap/quote`

```json
{
  "chain": "ethereum",
  "input_mint": "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
  "output_mint": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
  "amount": "1000000000000000000",
  "slippage": 100
}
```

(ETH → USDC, 1 ETH, 1% slippage)

**What to verify:** Same as Step 6

Related files:
- 1inch client: `backend/app/clients/oneinch.go`

---

## Step 8 — Auth (Get/Create User)

**Endpoint:** `GET /api/auth/me`
**Auth:** Firebase Bearer Token (required)

```
GET http://localhost:8080/api/auth/me
Authorization: Bearer <firebase_id_token>
```

**How to get a Firebase token:**
- Sign in via Google on the frontend, or
- Use Firebase Admin SDK to generate a custom token

**What to verify:**
- `data.id` — user ID in the database
- `data.email` — email from Google account
- `data.name` — name from Google account
- `data.avatar` — profile picture URL
- On first login, the user is automatically created in the DB

Related files:
- Middleware: `backend/app/middleware/firebase.go`
- Handler: `backend/app/handlers/auth.go`
- Service: `backend/app/services/user.go`
- Repository: `backend/app/repositories/user.go`
- Entity: `backend/app/entities/user.go`

---

## Step 9 — Follow Wallet (Watchlist)

**Endpoint:** `POST /api/watchlist/follow`
**Auth:** Firebase Bearer Token (required)

```json
{
  "wallet_address": "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
  "chain": "ethereum",
  "recommendation": "conditional_follow",
  "conditions": ["min_liquidity", "min_mcap"],
  "email_notify": true
}
```

**What to verify:**
- Response `status` = `"success"`
- Wallet appears in the watchlist (Step 11)
- `email_notify: true` means the user will receive email alerts when this wallet trades

Related files:
- Request DTO: `backend/app/dto/requests/watchlist.go`
- Handler: `backend/app/handlers/watchlist.go`
- Service: `backend/app/services/watchlist.go`
- Repository: `backend/app/repositories/watchlist.go`
- Entity: `backend/app/entities/watchlist.go`

---

## Step 10 — Update Watchlist Conditions

**Endpoint:** `PUT /api/watchlist/:address`
**Auth:** Firebase Bearer Token (required)

```json
{
  "conditions": ["min_liquidity", "min_pair_age", "min_volume"],
  "email_notify": false
}
```

**What to verify:**
- Conditions updated to match what you sent
- `email_notify` changed accordingly

---

## Step 11 — List Watchlist

**Endpoint:** `GET /api/watchlist`
**Auth:** Firebase Bearer Token (required)

**What to verify:**
- Array contains wallets you followed
- Each item has: `wallet_address`, `chain`, `recommendation`, `conditions`, `email_notify`

---

## Step 12 — Unfollow Wallet

**Endpoint:** `DELETE /api/watchlist/:address`
**Auth:** Firebase Bearer Token (required)

**What to verify:**
- Response is success
- Wallet no longer appears in the watchlist

---

## Step 13 — WebSocket Notifications

**Endpoint:** `ws://localhost:8080/ws?user_id=<USER_ID>`

Use a WebSocket client (Bruno doesn't support WS — use [websocat](https://github.com/nickel-org/websocat) or Postman):

```bash
websocat ws://localhost:8080/ws?user_id=1
```

**What to verify:**
- Connection succeeds (doesn't disconnect immediately)
- When a followed wallet makes a new trade, you receive a message:
  ```json
  {
    "type": "wallet_trade",
    "payload": {
      "wallet_address": "0x...",
      "chain": "ethereum",
      "token_address": "0x...",
      "token_symbol": "PEPE",
      "direction": "in",
      "value": "1000000",
      "traded_at": "2026-04-07T...",
      "liquidity": 150000,
      "market_cap": 500000,
      "price_change_24h": 12.5,
      "ai_assessment": "⚠️ This token launched 2 hours ago with $150k liquidity. Moderate risk — liquidity is decent but the token is still very new."
    }
  }
  ```
- Verify `ai_assessment` is present and contains a risk emoji (✅, ⚠️, or 🔴) with a beginner-friendly explanation
- Also check the database: notification is saved in the `notifications` table (including `ai_assessment` column)
- Monitor polls every 30 seconds (`backend/app/services/monitor.go`)

Related files:
- WebSocket handler: `backend/app/ws/handler.go`
- WebSocket hub: `backend/app/ws/hub.go`
- Monitor service: `backend/app/services/monitor.go`
- Monitor helper: `backend/app/services/monitor_helper.go`

---

## Step 14 — Email Notifications

If `email_notify: true` in watchlist and `RESEND_API_KEY` is set:

**What to verify:**
- When a followed wallet trades, check the user's email inbox
- Email contains: wallet address, token symbol, direction (bought/sold), amount, liquidity, market cap
- Subject format: `🔔 0xd8dA...96045 bought PEPE on ethereum`
- Email is sent async (does not block in-app notifications)
- If `RESEND_API_KEY` is still a placeholder, email will fail but in-app notifications still work (check log: `Monitor: email failed for user ...`)
- If `GEMINI_API_KEY` is not set, `ai_assessment` will be empty but notifications still send normally

Related files:
- Resend client: `backend/app/clients/resend.go`
- Interface: `backend/app/interfaces/notification.go`
- Wiring: `backend/router/container.go`
- Trigger: `backend/app/services/monitor_helper.go` (notifyFollowers method)

---

## Recommended Testing Order

```
1. Health Check          → make sure server is running
2. Analyze EVM           → test core feature, check all scoring metrics
3. Analyze Solana        → verify multi-chain works
4. Get Stored Wallet     → verify data is persisted in DB
5. Regenerate Insight    → test AI with different tones
6. Swap Quote Solana     → test Jupiter integration
7. Swap Quote EVM        → test 1inch integration
8. Auth (Get Me)         → test Firebase auth flow
9. Follow Wallet         → test watchlist creation
10. List Watchlist       → verify follow was successful
11. Update Conditions    → test updating conditions + email toggle
12. WebSocket Connect    → test real-time notifications
13. Wait for Monitor     → wait 30 seconds, check if notification arrives
14. Check Email          → if email_notify is true, check inbox
15. Unfollow Wallet      → test cleanup
```

---

## Tips

- Use a wallet address that is actively trading to test notifications (monitor needs new trades)
- If `ai_insight` is empty, make sure `GEMINI_API_KEY` is set with a valid key
- If all scores are 0, the wallet likely has no token transfers (only native ETH/SOL)
- Invalid limits automatically fall back to defaults (EVM: 10, Solana: 20)
- Bruno collection is available at `backend/api-docs/` for all endpoints
