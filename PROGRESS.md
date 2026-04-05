# Miora AI — Backend Progress

## ✅ Done

- Project structure (clean architecture: handlers, services, repositories, interfaces, entities, dto)
- Config loader (env-based, shared credentials antara Docker & Go app)
- Database setup (PostgreSQL via Docker Compose, GORM auto-migrate)
- Entities: Wallet, Transaction, WalletMetric
- Alchemy client: EVM (`alchemy_getAssetTransfers`) & Solana (`getSignaturesForAddress`)
- Repository layer (CRUD wallet, transactions, metrics)
- Service layer (fetch → save → score)
- Scoring system (formula sesuai README, masih placeholder)
- Handler + route: `POST /api/wallets/analyze`
- Health check: `GET /api/health`
- Error handling (AppError + ErrorResponse)
- Dockerfile (multi-stage build)
- Docker Compose (PostgreSQL)
- Makefile (git-commit, run-fe, run-be, run-all)
- .gitignore (Node, Next.js, Go, Foundry, Anchor, IDE)

## 🔲 Todo

- [ ] Real scoring logic (analisis PnL, win rate, risk dari data transaksi aktual)
- [ ] Middleware (rate limiting, auth/API key)
- [ ] WebSocket layer (real-time wallet tracking)
- [ ] GET endpoint: `/api/wallets/:address` (ambil hasil analisis yang sudah tersimpan)
- [ ] Pagination untuk transaction history
- [ ] AI layer integration (LLM untuk generate insight natural language)
- [ ] Caching (Redis untuk response yang sering diakses)
- [ ] DEX integration (swap via Jupiter/1inch)
- [ ] Smart alerts system
- [ ] Logging & monitoring
- [ ] Tests
