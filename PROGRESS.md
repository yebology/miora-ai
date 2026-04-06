# Miora AI — Backend Progress

## ✅ Done

- Project structure (clean architecture: handlers, services, repositories, interfaces, entities, dto)
- Config loader (env-based, shared credentials, no fallbacks, configurable scoring thresholds)
- Database setup (PostgreSQL via Docker Compose, GORM auto-migrate)
- Entities: Wallet, Transaction (with Direction + ContractAddress), WalletMetric
- Alchemy client: EVM (incoming + outgoing transfers) & Solana (getSignaturesForAddress)
- DexScreener client (token liquidity, market cap, pair age, price change)
- Moralis client (EVM historical token price by block + Solana current price)
- Birdeye client (Solana historical token price by unix timestamp)
- Repository layer (CRUD wallet, transactions, metrics)
- Service layer (fetch → enrich → FIFO buy-sell matching → PnL → score)
- Real scoring system (win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure)
- Scoring thresholds configurable from .env
- Risk exposure: informational only, not in final score formula
- Handler + route: `POST /api/wallets/analyze`
- Health check: `GET /api/health`
- Error handling (AppError constructors + output envelope)
- Validation (go-playground/validator + ParseAndValidateBody helper)
- DI container + router pattern (container.go + routes.go)
- Interfaces with `I` prefix (IWalletService, IWalletRepository)
- Dockerfile (multi-stage build)
- Docker Compose (PostgreSQL)
- Makefile (git-commit, run-fe, run-be, run-all, db-reset, db-seed)
- Migrations (auto-migrate, reset, seed)
- .gitignore (Node, Next.js, Go, Foundry, Anchor, IDE)
- Documentation (comments on all files)

## 🔲 Todo (Hackathon Priority)

- [ ] GET endpoint: `/api/wallets/:address` (retrieve stored analysis)
- [ ] AI layer (LLM for natural language insights + behavior classification)
- [ ] Tests (minimal)

## 📋 Post-Hackathon

- Middleware (rate limiting, auth/API key)
- WebSocket (real-time wallet tracking)
- Pagination for transaction history
- Caching (Redis)
- DEX integration (swap via Jupiter/1inch)
- Smart alerts system
- Logging & monitoring
