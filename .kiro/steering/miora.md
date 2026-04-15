# Miora AI — Project Steering

## Project Overview

Miora AI is a **Trading Reputation Protocol + AI Agent for Base**. It analyzes wallet trading behavior on Base, publishes scores on-chain via EAS attestation, and runs an AI agent that trades autonomously based on top-scored wallets.

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

### Two New Layers (V2)
- **EAS** — On-chain reputation attestations on Base Sepolia
- **AgentKit** — Autonomous AI trading agent via Coinbase Agentic Wallets

### Tech Stack
- **Backend**: Go + Fiber + GORM + WebSocket (Base only)
- **Agent Sidecar**: Python + FastAPI + Coinbase AgentKit
- **Frontend**: Next.js 16 + Tailwind CSS v4 + shadcn/ui + TypeScript
- **Database**: PostgreSQL via Docker Compose
- **Auth**: Firebase Auth (Google login)
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
| GET | `/api/auth/me` | Firebase | Get/create user |
| POST | `/api/watchlist/follow` | Firebase | Follow wallet |
| PUT | `/api/watchlist/:address` | Firebase | Update conditions |
| DELETE | `/api/watchlist/:address` | Firebase | Unfollow wallet |
| GET | `/api/watchlist` | Firebase | List followed wallets |
| GET | `/api/agent/status` | Firebase | Get agent status |
| PUT | `/api/agent/config` | Firebase | Update agent config |
| POST | `/api/agent/start` | Firebase | Start agent |
| POST | `/api/agent/pause` | Firebase | Pause agent |
| GET | `/api/agent/trades` | Firebase | Get agent trade history |
| WS | `/ws?user_id=ID` | — | Real-time notifications |

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

### Frontend (TypeScript/React)
- Next.js App Router
- Components by domain: `analyze/`, `watchlist/`, `landing/`, `layout/`, `ui/`, `agent/`
- shadcn/ui + `cn()` utility
- Dark mode default via next-themes
- API client in `lib/api.ts`
- Currently dummy data — replace with real API calls last

### Agent Sidecar (Python)
- FastAPI service wrapping `coinbase-agentkit`
- Runs on port 8090, called by Go backend via HTTP
- Handles: wallet creation, swap execution, balance queries

## Current Priority (Hackathon — Due April 27)

### 🔴 Must Have
1. **Deploy EAS attestation on Base Sepolia** — register schema, test attestation
2. **Build V2 frontend UI** — Agent page, reputation display (dummy data first)
3. **Record 1-minute founder video**
4. **Submit Base Batches application**

### 🟡 Should Have
5. **AgentKit proof of concept** — agent detects trade → executes on testnet
6. **Connect frontend to backend API** — replace dummy data (last)

### 🟢 Nice to Have
7. Agent dashboard with trade history
8. Reputation leaderboard

## What NOT to Do

- Do NOT add multi-chain support — Base only
- Do NOT add swap/DEX aggregator — agent handles trading
- Do NOT add tests unless explicitly asked
- Do NOT change scoring engine (`services/scoring.go`) — finalized
- Do NOT introduce new dependencies without checking existing ones
- Do NOT use `interface{}` — use `any` (Go 1.18+)
- Do NOT deploy to mainnet — Base Sepolia only
- Do NOT add custom smart contracts — use existing EAS contracts
