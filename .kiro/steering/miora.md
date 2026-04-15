# Miora AI — Project Steering

## Project Overview

Miora AI is a **Trading Reputation Protocol + AI Agent for Base**. It analyzes wallet trading behavior, publishes scores on-chain via EAS attestation, and runs an AI agent that trades autonomously based on top-scored wallets.

Target hackathon: **Base Batches 003 Student Track** (deadline April 27, 2026).

### Why This Concept (Research-Backed, Not Random)

Every decision in this project is backed by research documented in #[[file:MIORA_RESEARCH.md]]. Key findings:

1. **Why pivot from V1 (wallet analyzer) to V2 (reputation protocol)**: V1 was an iterative improvement on existing tools (Nansen, Arkham, Cielo). V2 creates a new primitive — on-chain trading reputation — that other protocols can build on.

2. **Why Base, not Solana**: Miora is DNA-EVM (5/6 chains are EVM including Base). Base Batches is application-based ("doesn't matter if code existed before"). Base 2026 strategy aligns: "discover what's trending from top traders."

3. **Why EAS + AgentKit + x402**: These are Base's own infrastructure (Coinbase-built). Using all three in one product = deep ecosystem alignment. No other BB003 project combines all three.

4. **Why trading reputation (not general reputation)**: Cred Protocol, ChainAware, Nomis focus on credit/fraud/general scoring. Nobody owns "trading quality reputation" on Base. Miora's FIFO PnL matching + 6-factor scoring = proprietary moat.

5. **Gap in BB003 cohort**: 12 selected teams cover AI agent infra, lending, privacy, prediction markets, neobank, FX. Nobody covers trading intelligence or wallet reputation. Miora fills this gap.

## Key Documents

- #[[file:README.md]] — Project README (V2 narrative)
- #[[file:PROGRESS.md]] — Development progress tracker (detailed per backend/frontend/smart contract)
- #[[file:MIORA_V2_CONCEPT.md]] — V2 concept document (pivot rationale, 3 layers, competitive analysis)
- #[[file:MIORA_RESEARCH.md]] — Research sources, step-by-step riset, competitive findings
- #[[file:BASE_BATCHES_APPLICATION.md]] — Base Batches application draft

## Architecture

### Three Layers
1. **Layer 1: On-chain Trading Reputation (EAS)** — Analyze wallet → score → publish attestation on Base Sepolia → queryable by any protocol
2. **Layer 2: Smart Follow + AI Alerts (Existing)** — Follow top wallets → real-time notifications with AI risk assessment → email alerts
3. **Layer 3: AI Trading Agent (AgentKit)** — Autonomous agent monitors top wallets → evaluates trades → executes swaps on Base via Agentic Wallet

### Tech Stack
- **Backend**: Go + Fiber + GORM + WebSocket. Clean architecture: handlers → services → repositories → interfaces.
- **Frontend**: Next.js 16 + Tailwind CSS v4 + shadcn/ui + TypeScript. Currently uses dummy data — needs to connect to real backend API.
- **Database**: PostgreSQL via Docker Compose.
- **Auth**: Firebase Auth (Google login).
- **AI**: Google Gemini (gemini-2.0-flash) for wallet insights and trade risk assessment.
- **On-chain**: EAS attestation on Base, Coinbase AgentKit for autonomous trading, x402 micropayments.
- **Smart Contracts**: Foundry (contracts/evm/) — EAS schema registration, optional helper contracts.

### API Endpoints (Existing)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/wallets/analyze` | Public | Analyze wallet address |
| GET | `/api/wallets/:address` | Public | Get stored analysis |
| POST | `/api/wallets/regenerate-insight` | Public | Regenerate AI insight |
| POST | `/api/swap/quote` | Public | Get swap quote (1inch) |
| GET | `/api/auth/me` | Firebase | Get/create user |
| POST | `/api/watchlist/follow` | Firebase | Follow wallet |
| PUT | `/api/watchlist/:address` | Firebase | Update conditions |
| DELETE | `/api/watchlist/:address` | Firebase | Unfollow wallet |
| GET | `/api/watchlist` | Firebase | List followed wallets |
| WS | `/ws?user_id=ID` | — | Real-time notifications |

### API Endpoints (V2 — To Build)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/reputation/:address` | Public | Get on-chain attestation data |
| GET | `/api/reputation/query` | x402 | Query reputation (micropayment) |
| POST | `/api/agent/start` | Firebase | Start AI trading agent |
| PUT | `/api/agent/config` | Firebase | Update agent config |
| POST | `/api/agent/pause` | Firebase | Pause/resume agent |
| GET | `/api/agent/status` | Firebase | Get agent status + trades |

## Coding Conventions

### Backend (Go)
- Follow existing clean architecture pattern: `clients/` → `dto/` → `entities/` → `handlers/` → `http/` → `interfaces/` → `services/` → `repositories/`
- All interfaces use `I` prefix (e.g., `IWalletService`, `IWalletRepository`)
- Services depend on repository interfaces, not concrete structs
- Handlers never contain business logic — only parse, delegate, and respond
- Use `pkg.AppError` for structured error handling
- Use `output.GetSuccess()` and `output.GetError()` for standardized API responses
- Use `utils.ParseAndValidateBody()` for request validation (go-playground/validator)
- All files must have package-level doc comments explaining purpose
- Config loaded from `.env` via `godotenv` — no fallback values, all required keys must be set
- New V2 files follow same pattern: create interface → implement client/service → create handler → register routes → wire in container

### Frontend (TypeScript/React)
- Next.js App Router (not Pages Router)
- Components organized by domain: `analyze/`, `watchlist/`, `landing/`, `layout/`, `ui/`, `agent/` (V2 new)
- Use shadcn/ui components from `components/ui/`
- Use `cn()` utility from `lib/utils.ts` for conditional classnames
- Dark mode default via next-themes
- Types in `types/` directory matching backend response shapes
- API client in `lib/api.ts` — all backend calls go through this
- Currently dummy data in `constants/dummy.ts` and `constants/dummy-watchlist.ts` — replace with real API calls
- Use `@tanstack/react-query` for data fetching (already installed, not yet used)

### Smart Contracts (Foundry)
- EVM contracts in `contracts/evm/`
- Use Foundry for build, test, deploy
- EAS attestations use existing EAS contracts on Base — no custom contract needed for basic attestation
- Optional helper contracts for batch queries or agent guardrails
- Deploy to Base Sepolia for hackathon demo

## Current Priority (Hackathon — Due April 27)

### 🔴 Must Have (in order)
1. **Deploy EAS attestation on Base Sepolia** — this is the V2 differentiator, proof of on-chain reputation
   - Register EAS schema → create `clients/eas.go` → integrate into scoring flow → show attestation in frontend
2. **Build V2 frontend UI** — Agent page, reputation display (using dummy data first so UI can be reviewed visually)
3. **Record 1-minute founder video**
4. **Submit Base Batches application**

### 🟡 Should Have
5. **Basic AgentKit proof of concept** — detect trade → evaluate risk → execute swap on testnet
6. **x402 reputation query endpoint** — monetization proof
7. **Connect frontend to backend API** — replace dummy data with real API calls (do this last, after UI is finalized)

### 🟢 Nice to Have
8. Agent dashboard with trade history
9. Reputation leaderboard page
10. Multiple wallet monitoring in agent

## What's Already Done (Don't Touch)
- ✅ Solana/V1 cleanup — all Solana code removed, backend compiles clean, frontend builds clean
- ✅ Scoring engine (`services/scoring.go`) — 6-factor scoring, FIFO PnL, finalized
- ✅ AI layer (`services/ai.go`, `clients/gemini.go`) — Gemini insights + trade risk assessment
- ✅ Wallet monitoring (`services/monitor.go`) — background polling, condition checking, notification dispatch
- ✅ Watchlist system — follow/unfollow, conditions, CRUD
- ✅ Auth system — Firebase middleware, user service
- ✅ All frontend pages and components — analyze, watchlist, landing, layout, providers, UI
- ✅ API client (`lib/api.ts`) — all endpoint functions defined, ready to use
- ✅ Landing page — V2 narrative applied

## What NOT to Do

- Do NOT add Solana-specific code — project has pivoted to Base-only
- Do NOT add tests unless explicitly asked
- Do NOT change the scoring engine logic (`services/scoring.go`) — it's finalized
- Do NOT introduce new dependencies without checking if existing ones cover the use case
- Do NOT use `interface{}` — use `any` instead (Go 1.18+)
- Do NOT rewrite existing working code — extend it for V2 features
- Do NOT deploy to mainnet — use Base Sepolia for all V2 features
