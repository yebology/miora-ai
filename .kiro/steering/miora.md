# Miora AI — Project Steering

## Project Overview

Miora AI is a **Trading Reputation Protocol + AI Bots for Base**. It analyzes wallet trading behavior on Base, publishes scores on-chain via EAS attestation, and runs per-wallet AI bots that trade autonomously based on top-scored wallets.

Target hackathon: **Base Batches 003 Student Track** (deadline April 27, 2026).
Chain: **Base only** (Base Sepolia for testnet).

## Key Documents

- #[[file:README.md]] — Project README
- #[[file:PROGRESS.md]] — Development progress tracker
- #[[file:MIORA_V2_CONCEPT.md]] — V2 concept document
- #[[file:MIORA_RESEARCH.md]] — Research sources
- #[[file:BASE_BATCHES_APPLICATION.md]] — Base Batches application draft

## Architecture

### Three Layers
1. **Layer 1: On-chain Trading Reputation (EAS)** — Analyze wallet → score → publish attestation on Base Sepolia
2. **Layer 2: Smart Follow + AI Alerts** — Follow top wallets → real-time notifications with AI risk assessment
3. **Layer 3: AI Trading Bots (AgentKit)** — Per-wallet bots that copy trades (buy + sell) or consensus bots that trade on multi-wallet agreement

### Bot Types
- **Wallet Bot** — Targets one specific wallet. Copies its buys AND sells. User selects wallet from watchlist, conditions auto-filled from analyze result.
- **Consensus Bot** — Scans all Miora-analyzed wallets. Trades when multiple high-score wallets buy the same token within a configurable time window. Has its own budget, min_score, threshold, and time window.

### Auth
- **Wallet-based auth** — MetaMask connect via wagmi/viem
- **X-Wallet-Address header** — sent with all protected API requests
- **No Firebase** — no Google login, no Firebase SDK

### Tech Stack
- **Backend**: Go + Fiber + GORM + WebSocket (Base only)
- **Agent Sidecar**: Python + FastAPI + Coinbase AgentKit
- **Frontend**: Next.js 16 + Tailwind CSS v4 + shadcn/ui + TypeScript
- **Database**: PostgreSQL via Docker Compose
- **Auth**: Wallet-based (MetaMask connect, X-Wallet-Address header)
- **AI**: Google Gemini (gemini-2.0-flash)
- **On-chain**: EAS attestation, Coinbase AgentKit + Agentic Wallets
- **Data**: Alchemy (Base RPC), DexScreener (pair data), Moralis (historical prices)

### API Endpoints
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/wallets/analyze` | Public | Analyze wallet on Base |
| GET | `/api/wallets/:address` | Public | Get stored analysis |
| POST | `/api/wallets/regenerate-insight` | Public | Regenerate AI insight |
| GET | `/api/reputation/:address` | Public | Get on-chain attestation |
| GET | `/api/auth/me` | Wallet | Get/create user |
| POST | `/api/watchlist/follow` | Wallet | Follow wallet |
| PUT | `/api/watchlist/:address` | Wallet | Update conditions |
| DELETE | `/api/watchlist/:address` | Wallet | Unfollow wallet |
| GET | `/api/watchlist` | Wallet | List followed wallets |
| POST | `/api/agent/bots` | Wallet | Create a new bot (wallet or consensus type) |
| GET | `/api/agent/bots` | Wallet | List all user's bots |
| GET | `/api/agent/bots/:id` | Wallet | Get bot details + status |
| PUT | `/api/agent/bots/:id` | Wallet | Update bot configuration |
| POST | `/api/agent/bots/:id/start` | Wallet | Start a bot |
| POST | `/api/agent/bots/:id/pause` | Wallet | Pause a bot |
| DELETE | `/api/agent/bots/:id` | Wallet | Delete a bot |
| GET | `/api/agent/bots/:id/trades` | Wallet | Get bot's trade history |
| WS | `/ws?wallet_address=ADDR` | — | Real-time notifications |

## Coding Conventions

### Backend (Go)
- Clean architecture: `clients/` → `dto/` → `entities/` → `handlers/` → `http/` → `interfaces/` → `services/` → `repositories/`
- All interfaces use `I` prefix (e.g., `IWalletService`)
- Services depend on repository interfaces, not concrete structs
- Handlers never contain business logic — only parse, delegate, respond
- Use `pkg.AppError` for structured error handling
- Use `output.GetSuccess()` and `output.GetError()` for API responses
- Config loaded from `.env` via `godotenv`
- Base only — no multi-chain logic
- Auth via `middleware/wallet_auth.go` — reads `X-Wallet-Address` header
- No Firebase SDK — no Firebase imports anywhere

### Frontend (TypeScript/React)
- Next.js App Router
- Components by domain: `analyze/`, `watchlist/`, `landing/`, `layout/`, `ui/`, `agent/`
- shadcn/ui + `cn()` utility
- Dark mode default via next-themes
- API client in `lib/api.ts` — sends `X-Wallet-Address` header
- Auth via wallet connect (MetaMask, wagmi/viem) — no Firebase, no Google login
- Currently dummy data — replace with real API calls last

### Agent Sidecar (Python)
- FastAPI service wrapping `coinbase-agentkit`
- Runs on port 8090, called by Go backend via HTTP
- Handles: wallet creation, swap execution, balance queries

### Bot Config Model
- `bot_type`: `"wallet"` or `"consensus"`
- Wallet bot fields: `target_wallet`, `budget`, `max_per_trade`, `min_score`, `conditions`
- Consensus bot fields: `budget`, `max_per_trade`, `min_score`, `consensus_threshold`, `consensus_window_min`
- No `risk_tolerance` field — conditions from analyze result handle risk filtering
- Each bot is independent — user can have multiple bots

## Current Priority (Hackathon — Due April 27)

### 🔴 Must Have
1. **Deploy EAS attestation on Base Sepolia** — register schema, test attestation
2. **Build V2 frontend UI** — Bot page, reputation display (dummy data first)
3. **Record 1-minute founder video**
4. **Submit Base Batches application**

### 🟡 Should Have
5. **AgentKit proof of concept** — bot detects trade → executes on testnet
6. **Connect frontend to backend API** — replace dummy data (last)

### 🟢 Nice to Have
7. Bot dashboard with trade history
8. Reputation leaderboard

## What NOT to Do

- Do NOT add multi-chain support — Base only
- Do NOT add swap/DEX aggregator — bots handle trading
- Do NOT add tests unless explicitly asked
- Do NOT change scoring engine (`services/scoring.go`) — finalized
- Do NOT introduce new dependencies without checking existing ones
- Do NOT use `interface{}` — use `any` (Go 1.18+)
- Do NOT deploy to mainnet — Base Sepolia only
- Do NOT add custom smart contracts — use existing EAS contracts
- Do NOT add Firebase — auth is wallet-based only
- Do NOT add risk_tolerance to bot config — conditions handle risk filtering
