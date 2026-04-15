# Miora AI — Progress Tracker

## 🔄 Pivot: V1 → V2

Miora has pivoted from a multi-chain wallet analyzer + DEX aggregator (V1) to a **Trading Reputation Protocol on Base** (V2).

- **V1**: Analytics tool — analyze wallet, follow, get alerts, swap manually (off-chain only)
- **V2**: Infrastructure protocol — publish trading scores on-chain via EAS, autonomous AI agent trading via AgentKit, monetize reputation API via x402

### Why Pivot
- V1 was an iterative improvement on existing tools (Nansen, Arkham, Cielo)
- V2 creates a new primitive: on-chain trading reputation that other protocols can build on
- V2 deeply integrates with Base ecosystem (EAS, AgentKit, x402)
- V2 aligns with Base 2026 strategy: global markets, stablecoin payments, home for builders/agents

---

## 🏗️ BACKEND (Go + Fiber + GORM)

### ✅ Core Scoring Engine — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `services/scoring.go` | ✅ Done | Multi-factor scoring (0-100): win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure |
| `services/wallet.go` | ✅ Done | Full analysis orchestration: fetch txs → enrich data → calculate PnL → score → recommend → AI insight |
| `services/wallet_helper.go` | ✅ Done | FIFO buy-sell matching for PnL, token data fetching, condition generation, price helpers |
| — | ✅ Done | 3-tier recommendations: `full_follow` (80-100), `conditional_follow` (40-79), `avoid` (<40) |
| — | ✅ Done | Dynamic conditional thresholds computed from wallet's own token data (median liquidity, mcap, volume, avg pair age) |

### ✅ AI Layer — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `services/ai.go` | ✅ Done | Gemini-powered insight generation (simple, eli5, custom tone) + trade risk assessment |
| `clients/gemini.go` | ✅ Done | Google Gemini API client (gemini-2.0-flash) |
| `dto/prompts/wallet.go` | ✅ Done | Prompt templates for wallet insights and trade assessments |

### ✅ Wallet Monitoring & Alerts — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `services/monitor.go` | ✅ Done | Background polling of watched wallets, new trade detection |
| `services/monitor_helper.go` | ✅ Done | Condition checking, notification dispatch logic |
| `ws/hub.go` | ✅ Done | WebSocket connection hub (broadcast trade notifications) |
| `ws/handler.go` | ✅ Done | WebSocket upgrade and connection handlers |
| `clients/resend.go` | ✅ Done | Email notifications via Resend (async, non-blocking) |
| — | ✅ Done | AI risk assessment per trade notification (Gemini evaluates token before notifying) |

### ✅ Watchlist System — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `services/watchlist.go` | ✅ Done | Follow/unfollow, update conditions, list watchlist |
| `handlers/watchlist.go` | ✅ Done | CRUD endpoints: POST /follow, DELETE /:address, GET /, PUT /:address |
| `repositories/watchlist.go` | ✅ Done | Create, Delete, Exists, FindByUser, Update |
| `entities/watchlist.go` | ✅ Done | UserID, WalletAddress, Chain, Recommendation, Conditions (JSON), EmailNotify |

### ✅ Auth & User System — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `services/user.go` | ✅ Done | Find-or-create from Firebase UID |
| `handlers/auth.go` | ✅ Done | GET /auth/me endpoint |
| `middleware/firebase.go` | ✅ Done | Firebase token verification middleware |
| `repositories/user.go` | ✅ Done | FindByFirebaseUID, Create, Update |
| `entities/user.go` | ✅ Done | ID, FirebaseUID, Email, Name, Avatar |

### ✅ Swap System — REMOVED
> Swap system (1inch, manual trading) has been removed. Agent handles all trading via AgentKit.

### ✅ Data Clients — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `clients/evm.go` | ✅ Done | Alchemy EVM client (Base only) |
| `clients/dexscreener.go` | ✅ Done | Token pair data (liquidity, mcap, pair age, price change) |
| `clients/moralis.go` | ✅ Done | Historical token prices by block on Base |
| `clients/helpers.go` | ✅ Done | Shared HTTP client helpers |

### ✅ Entities (Database Models) — DONE, Needs V2 Additions
| File | Status | Description |
|------|--------|-------------|
| `entities/user.go` | ✅ Done | ID, FirebaseUID, Email, Name, Avatar |
| `entities/wallet.go` | ✅ Done | ID, Address, Chain |
| `entities/transaction.go` | ✅ Done | ID, WalletID, Hash, Chain, From, To, Value, TokenSymbol, ContractAddress, Direction, BlockNumber, Timestamp |
| `entities/wallet_metric.go` | ✅ Done | ID, WalletID, TotalTransactions, 6 scoring metrics, FinalScore, Recommendation |
| `entities/watchlist.go` | ✅ Done | ID, UserID, WalletAddress, Chain, Recommendation, Conditions (JSON), EmailNotify |
| `entities/notification.go` | ✅ Done | ID, UserID, WalletAddress, Chain, TokenAddress, TokenSymbol, Direction, Value, Liquidity, MarketCap, AiAssessment, Read |

### ✅ Infrastructure — DONE, No Changes Needed
| File | Status | Description |
|------|--------|-------------|
| `router/container.go` | ✅ Done | DI container: Clients → Repos → Services → Handlers |
| `router/routes.go` | ✅ Done | Route registration + middleware setup |
| `config/config.go` | ✅ Done | Environment config loader (all required keys) |
| `constants/chains.go` | ✅ Done | Base only (multi-chain removed) |
| `constants/limits.go` | ✅ Done | Configurable transaction limits per chain |
| `constants/error.go` | ✅ Done | Error message constants |
| `constants/success.go` | ✅ Done | Success message constants |
| `pkg/error.go` | ✅ Done | AppError struct for structured error handling |
| `utils/helper.go` | ✅ Done | General helpers (clamp, round, etc.) |
| `utils/math.go` | ✅ Done | Math utilities (stdDev, median) |
| `utils/utils.go` | ✅ Done | Validation, parsing, response helpers |
| `migrations/migrations.go` | ✅ Done | Auto-migrate all entities |
| `migrations/reset.go` | ✅ Done | Drop all tables and re-apply |
| `migrations/seed.go` | ✅ Done | Seed development data |
| `Dockerfile` | ✅ Done | Docker image for backend |
| `docker-compose.yml` | ✅ Done | PostgreSQL + backend services |
| `main.go` | ✅ Done | Entry point: config → DB → migrations → Fiber → routes |

### ✅ Solana/V1 Cleanup — DONE
- [x] Removed all Solana-specific code: `clients/solana.go`, `clients/birdeye.go`, `clients/jupiter.go`, `interfaces/birdeye.go`
- [x] Removed Solana references from config, constants, services, router, DTOs, migrations, seed
- [x] Backend compiles clean

### 🆕 V2 Backend — Layer 1: EAS Attestation (On-chain Reputation)
- [x] Research EAS SDK for Go (or use direct contract call via go-ethereum/abigen)
- [x] Add `AttestationUID` field to `entities/wallet_metric.go`
- [x] Add `SchemaUID` to config (env var: `EAS_SCHEMA_UID`)
- [x] Add `EAS_CONTRACT_ADDRESS` to config (Base Sepolia EAS contract)
- [x] Add `BASE_RPC_URL` to config (Base Sepolia RPC endpoint)
- [x] Add `ATTESTER_PRIVATE_KEY` to config (wallet private key for signing attestations)
- [x] Create `clients/eas.go` — EAS client: create attestation, query attestation by UID
- [x] Create `interfaces/eas.go` — `IEASClient` interface
- [x] Update `services/wallet.go` — after scoring, call EAS client to publish attestation on-chain
- [x] Create `handlers/reputation.go` — GET /reputation/:address (return attestation data + on-chain proof)
- [x] Create `http/reputation.go` — register reputation routes (public)
- [x] Create `dto/responses/reputation.go` — reputation response DTO (score, recommendation, attestation UID, txn hash, timestamp)
- [x] Wire EAS client into `router/container.go`
- [x] Register reputation routes in `router/routes.go`
- [x] Update `migrations/migrations.go` to auto-migrate updated WalletMetric
- [ ] Register EAS schema on Base Sepolia → run `make register-schema` (needs testnet ETH in attester wallet)
- [ ] Set `EAS_SCHEMA_UID` in `.env` (printed by register-schema command)
- [ ] Test end-to-end: analyze wallet → attestation published → verify on [BaseScan](https://base-sepolia.easscan.org)

### 🆕 V2 Backend — Layer 2: x402 Reputation API (Monetization)
- [x] Research x402 protocol integration for Go (using `mark3labs/x402-go` library)
- [x] Create `middleware/x402.go` — x402 payment verification middleware (Fiber-compatible)
- [x] Add `X402_RECIPIENT_ADDRESS` to config (USDC receiving address)
- [x] Add `X402_PRICE_USDC` to config (price per query in USDC)
- [x] Create `handlers/reputation.go` → `QueryReputation()` — GET /reputation/query?address=0x... (x402-protected)
- [x] Update `http/reputation.go` — register x402-protected reputation routes
- [x] Create `dto/responses/reputation_query.go` — query response DTO
- [x] Wire x402 middleware into `router/routes.go`
- [ ] Add `X402_RECIPIENT_ADDRESS` to `backend/.env` (actual wallet address)
- [ ] Test with x402 client on Base Sepolia

### 🆕 V2 Backend — Layer 3: AI Trading Agent (AgentKit)
- [x] Research Coinbase AgentKit SDK — only available in TypeScript/Python, no Go SDK
- [x] Create `entities/agent_config.go` — AgentConfig entity: UserID, Budget, MaxPerTrade, RiskTolerance, MinScore, Conditions (JSON), Status (active/paused/stopped), AgentWalletAddress, TotalSpent, TotalTrades
- [x] Create `entities/agent_trade.go` — AgentTrade entity: AgentConfigID, SourceWallet, SourceScore, TokenAddress, TokenSymbol, Direction, AmountUSD, TxHash, Status, Reason, RiskAssessment
- [x] Create `repositories/agent.go` — AgentConfig + AgentTrade CRUD
- [x] Create `interfaces/agent.go` — `IAgentRepository`, `IAgentService`
- [x] Create `services/agent.go` — Agent service: GetOrCreateConfig, UpdateConfig, Start, Pause, GetStatus, GetTrades
- [x] Create `handlers/agent.go` — GET /agent/status, PUT /agent/config, POST /agent/start, POST /agent/pause, GET /agent/trades
- [x] Create `http/agent.go` — register agent routes (protected, Firebase auth)
- [x] Create `dto/requests/agent.go` — agent config request DTO
- [x] Wire agent service + handler into `router/container.go`
- [x] Register agent routes in `router/routes.go`
- [x] Update `migrations/migrations.go` to auto-migrate AgentConfig + AgentTrade
- [x] Create Python AgentKit sidecar (`agent/main.py`) — FastAPI service wrapping `coinbase-agentkit`
- [x] Create `clients/agentkit.go` — Go HTTP client to call Python sidecar (GetWallet, ExecuteSwap, IsHealthy)
- [x] Create `services/agent_loop.go` — Background loop: poll active configs → check top wallets → evaluate conditions → execute swap via sidecar
- [x] Add swap endpoint to Python sidecar (`POST /swap`)
- [x] Add `make setup-agent` and `make run-agent` to Makefile
- [ ] Add `CDP_API_KEY_ID` and `CDP_API_KEY_SECRET` to `agent/.env`
- [ ] Test end-to-end: agent detects trade → evaluates → executes swap on Base Sepolia

---

## 🎨 FRONTEND (Next.js 16 + Tailwind v4 + shadcn/ui)

### ✅ Tech Stack — DONE
- Next.js 16 + Tailwind CSS v4 + shadcn/ui + TypeScript
- next-themes (dark/light mode, default dark), Space Grotesk font
- wagmi + viem + @reown/appkit (wallet connect — MetaMask, WalletConnect)
- @tanstack/react-query (installed, not yet used — ready for API integration)
- @coinbase/cdp-sdk (NOT installed in package.json — needs to be added for AgentKit)

### ✅ Pages — DONE, Needs V2 Additions
| Page | Status | Description |
|------|--------|-------------|
| `app/page.tsx` (Landing) | ✅ Done | V2 hero + narrative (reputation protocol + AI agent) |
| `app/analyze/page.tsx` | ✅ Done | Wallet analysis page — currently uses dummy data |
| `app/watchlist/page.tsx` | ✅ Done | Watchlist dashboard — currently uses dummy data |
| `app/watchlist/[chain]/[address]/page.tsx` | ✅ Done | Watchlist detail page — currently uses dummy data |
| `app/swap/page.tsx` | ⚠️ To remove | Swap system removed — agent handles trading |
| `app/login/page.tsx` | ⚠️ Placeholder | Login page skeleton |
| `app/agent/page.tsx` | ❌ To build | Agent setup + dashboard page (V2 new) |

### ✅ Components — DONE, Needs V2 Additions
| Directory | Files | Status |
|-----------|-------|--------|
| `components/analyze/` | analyze-form, analysis-result, score-ring, metric-bar, ai-insight-card, conditions-card, traded-tokens-table, recommendation-badge | ✅ Done — all use dummy data, have TODO comments for real API |
| `components/watchlist/` | watchlist-card, notification-item | ✅ Done — uses dummy data |
| `components/landing/` | hero-section, hero-background, features-section, how-it-works-section, chains-section, cta-section | ✅ Done — V2 narrative applied |
| `components/layout/` | navbar, footer, theme-toggle | ✅ Done — needs "Agent" nav item |
| `components/providers/` | auth-provider, theme-provider, web3-provider | ✅ Done — auth-provider uses simulated login |
| `components/ui/` | button, card, badge, dialog, input, label, progress, select, sheet, auth-guard-modal, wallet-guard-modal | ✅ Done |
| `components/icons/` | google | ✅ Done |
| `components/agent/` | — | ❌ To build (V2 new) |

### ✅ Data Layer — DONE, Needs Real API Connection
| File | Status | Description |
|------|--------|-------------|
| `lib/api.ts` | ✅ Done | API client with endpoints (analyzeWallet, getWallet, regenerateInsight, getWatchlist, followWallet, unfollowWallet, updateWatchlist, getMe) |
| `constants/dummy.ts` | ⚠️ Dummy | 3 dummy wallet analyses (conditional_follow, full_follow, avoid) — to be replaced |
| `constants/dummy-watchlist.ts` | ⚠️ Dummy | 3 dummy watchlist items + 4 dummy notifications — to be replaced |
| `constants/landing.ts` | ✅ Done | Landing page copy (V2 narrative, Solana removed) |
| `constants/nav.ts` | ✅ Done | Navigation items — needs "Agent" added |
| `constants/tokens.ts` | ⚠️ To remove | Token list for swap — swap system removed |
| `types/wallet.ts` | ✅ Done | WalletAnalysis, TradedToken, Condition types |
| `types/watchlist.ts` | ✅ Done | WatchlistItem, Notification types |
| `types/api.ts` | ✅ Done | ApiResponse envelope type |
| `types/swap.ts` | ⚠️ To remove | SwapQuote type — swap system removed |
| `hooks/use-animate-on-scroll.ts` | ✅ Done | Scroll animation hook |
| `lib/utils.ts` | ✅ Done | cn() utility for conditional classnames |

### ✅ Solana/V1 Cleanup — DONE
- [x] Removed Solana from frontend: chain selectors, tokens, dummy data, landing page, swap page
- [x] Deleted `solana.svg`
- [x] Updated hero section + README tagline to V2 narrative
- [x] Frontend builds clean

### 🆕 V2 Frontend — Agent Page & Components
- [ ] Add "Agent" to `constants/nav.ts` navigation items
- [ ] Create `app/agent/page.tsx` — Agent setup + dashboard page
- [ ] Create `types/agent.ts` — AgentConfig, AgentTrade, AgentStatus types
- [ ] Add agent API functions to `lib/api.ts` — startAgent, pauseAgent, getAgentStatus, updateAgentConfig, getAgentTrades
- [ ] Create `components/agent/agent-config-form.tsx` — Budget, max per trade, risk tolerance, conditions form
- [ ] Create `components/agent/agent-status-card.tsx` — Agent status (active/paused/stopped), wallet balance, total trades
- [ ] Create `components/agent/agent-trade-history.tsx` — Table of agent's executed trades with PnL
- [ ] Create `components/agent/agent-wallet-card.tsx` — Agentic wallet address, balance, deposit/withdraw

### 🆕 V2 Frontend — Reputation Display
- [ ] Add attestation badge/link to `components/analyze/analysis-result.tsx` — show EAS attestation UID + BaseScan link after analysis
- [ ] Create `components/analyze/attestation-badge.tsx` — "Verified on Base" badge with attestation link

### 🔌 Connect Frontend to Backend API (Replace Dummy Data) — LAST PRIORITY
> Do this after all V2 UI is built and backend is running. Currently using dummy data so the UI can be reviewed visually first.

#### Auth Provider (`components/providers/auth-provider.tsx`)
- [ ] Replace simulated sign-in with real Firebase Google sign-in (`signInWithPopup`, `GoogleAuthProvider`)
- [ ] Store Firebase ID token for API calls
- [ ] Call `getMe()` after sign-in to sync user with backend
- [ ] Replace simulated sign-out with real Firebase sign-out
- [ ] Add `getToken()` method to auth context for components to use

#### Analyze Page (`app/analyze/page.tsx`)
- [ ] Replace dummy data simulation with real `analyzeWallet()` call from `lib/api.ts`
- [ ] Replace dummy "wallet exists" check with real `getWallet()` call
- [ ] Remove `DUMMY_ANALYSIS`, `DUMMY_FULL_FOLLOW`, `DUMMY_AVOID` imports
- [ ] Handle loading, error, and empty states from real API responses

#### AI Insight Card (`components/analyze/ai-insight-card.tsx`)
- [ ] Replace dummy insight regeneration with real `regenerateInsight()` call from `lib/api.ts`
- [ ] Remove hardcoded `dummyInsights` object

#### Conditions Card (`components/analyze/conditions-card.tsx`)
- [ ] Replace dummy follow action with real `followWallet()` call from `lib/api.ts`
- [ ] Pass Firebase auth token from `useAuth()` context

#### Analysis Result (`components/analyze/analysis-result.tsx`)
- [ ] Replace dummy "Follow Wallet" action with real `followWallet()` call
- [ ] Pass Firebase auth token from `useAuth()` context

#### Watchlist Page (`app/watchlist/page.tsx`)
- [ ] Replace `DUMMY_WATCHLIST` with real `getWatchlist()` call from `lib/api.ts`
- [ ] Replace `DUMMY_NOTIFICATIONS` with real notification data (WebSocket or polling)
- [ ] Replace dummy unfollow with real `unfollowWallet()` call
- [ ] Replace dummy toggle notify with real `updateWatchlist()` call
- [ ] Pass Firebase auth token from `useAuth()` context

#### Watchlist Detail Page (`app/watchlist/[chain]/[address]/page.tsx`)
- [ ] Replace `DUMMY_ANALYSIS` with real `getWallet()` call from `lib/api.ts`
- [ ] Replace `DUMMY_NOTIFICATIONS` with real notification data filtered by wallet
- [ ] Replace dummy re-analyze with real `analyzeWallet()` call
- [ ] Remove all dummy data imports

---

## 🧹 Cleanup — DONE
- [x] Removed swap system (1inch, swap handler, swap service, swap routes)
- [x] Removed multi-chain support (Base only — removed Ethereum, Arbitrum, Optimism, Polygon from chain registry)
- [x] Removed smart contract placeholders (Counter.sol, Counter.t.sol, Counter.s.sol)
- [x] Removed Solana/V1 code (all Solana clients, interfaces, references)
- [ ] Remove `frontend/constants/dummy.ts` after API connection is done
- [ ] Remove `frontend/constants/dummy-watchlist.ts` after API connection is done

---

## 🔧 Infrastructure & Config Updates
- [x] Add EAS env vars to `backend/config/config.go`
- [x] Add x402 env vars to `backend/config/config.go`
- [x] Add Agent env vars to `backend/config/config.go`
- [x] Update `Makefile` with `make register-schema`, `make setup-agent`, `make run-agent`
- [x] Remove `ONEINCH_API_KEY` from config (swap system removed)
- [x] Update `constants/chains.go` to Base only
- [ ] Add EAS env vars to `backend/.env` (actual values)
- [ ] Add `X402_RECIPIENT_ADDRESS` to `backend/.env`
- [ ] Add `CDP_API_KEY_ID` and `CDP_API_KEY_SECRET` to `agent/.env`
- [ ] Add `@coinbase/cdp-sdk` to `frontend/package.json` (if frontend needs AgentKit)
- [ ] Update `frontend/.env` with any new public env vars


