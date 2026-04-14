# Miora AI — Project Steering

## Project Overview

Miora AI is a **Trading Reputation Protocol + AI Agent for Base**. It analyzes wallet trading behavior, publishes scores on-chain via EAS attestation, and runs an AI agent that trades autonomously based on top-scored wallets.

Target hackathon: **Base Batches 003 Student Track** (deadline April 27, 2026).

### Why This Concept (Research-Backed, Not Random)

Every decision in this project is backed by research documented in #[[file:MIORA_RESEARCH.md]]. Key findings:

1. **Why pivot from V1 (wallet analyzer) to V2 (reputation protocol)**: V1 was an iterative improvement on existing tools (Nansen, Arkham, Cielo). Scored 6.5/10 for Solana Frontier. V2 creates a new primitive — on-chain trading reputation — that other protocols can build on.

2. **Why Base, not Solana**: Miora is DNA-EVM (5/6 chains are EVM including Base). Base fit scored 7.5/10 vs Solana 5.5/10. Base Batches is application-based ("doesn't matter if code existed before"), better for current state. Base 2026 strategy aligns: "discover what's trending from top traders."

3. **Why EAS + AgentKit + x402**: These are Base's own infrastructure (Coinbase-built). Using all three in one product = deep ecosystem alignment. No other BB003 project combines all three. Juri will see: "he uses our tech stack cohesively."

4. **Why trading reputation (not general reputation)**: Cred Protocol, ChainAware, Nomis focus on credit/fraud/general scoring. Nobody owns "trading quality reputation" on Base. Miora's FIFO PnL matching + 6-factor scoring = proprietary moat.

5. **Gap in BB003 cohort**: 12 selected teams cover AI agent infra, lending, privacy, prediction markets, neobank, FX. Nobody covers trading intelligence or wallet reputation. Miora fills this gap.

## Key Documents

- #[[file:README.md]] — Project README (V2 narrative)
- #[[file:PROGRESS.md]] — Development progress tracker
- #[[file:MIORA_V2_CONCEPT.md]] — V2 concept document (pivot rationale, 3 layers, competitive analysis)
- #[[file:MIORA_RESEARCH.md]] — Research sources, step-by-step riset, competitive findings
- #[[file:BASE_BATCHES_APPLICATION.md]] — Base Batches application draft
- #[[file:CLEANUP.md]] — Codebase cleanup checklist (remove Solana/V1 code)

## Architecture

- **Backend**: Go + Fiber + GORM + WebSocket. Clean architecture: handlers → services → repositories → interfaces.
- **Frontend**: Next.js 16 + Tailwind CSS v4 + shadcn/ui + TypeScript. Currently uses dummy data — needs to connect to real backend API.
- **Database**: PostgreSQL via Docker Compose.
- **Auth**: Firebase Auth (Google login).
- **AI**: Google Gemini (gemini-2.0-flash) for wallet insights and trade risk assessment.
- **On-chain (V2 new)**: EAS attestation on Base, Coinbase AgentKit for autonomous trading, x402 micropayments.

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

### Frontend (TypeScript/React)
- Next.js App Router (not Pages Router)
- Components organized by domain: `analyze/`, `watchlist/`, `landing/`, `layout/`, `ui/`
- Use shadcn/ui components from `components/ui/`
- Use `cn()` utility from `lib/utils.ts` for conditional classnames
- Dark mode default via next-themes
- Types in `types/` directory matching backend response shapes
- Currently dummy data in `constants/dummy.ts` and `constants/dummy-watchlist.ts` — replace with real API calls

### Smart Contracts
- EVM contracts use Foundry (contracts/evm/)
- Solana contracts removed (V2 pivot — Base only)

## Current Priority

1. **Connect frontend to backend API** (replace dummy data) — this is the #1 blocker
2. **Deploy EAS attestation on Base Sepolia** (proof of on-chain reputation)
3. **Basic AgentKit proof of concept**
4. **Record founder video**
5. **Submit Base Batches application by April 27**

## What NOT to Do

- Do NOT add Solana-specific code — project has pivoted to Base-only
- Do NOT add tests unless explicitly asked
- Do NOT change the scoring engine logic (services/scoring.go) — it's finalized
- Do NOT introduce new dependencies without checking if existing ones cover the use case
- Do NOT use `interface{}` — use `any` instead (Go 1.18+)
