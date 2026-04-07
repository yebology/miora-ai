# Miora AI — Progress

## ✅ Backend Done

- Clean architecture (handlers, services, repositories, interfaces, entities, dto)
- Config loader (env-based, no fallbacks, configurable scoring thresholds)
- Database (PostgreSQL via Docker Compose, GORM auto-migrate)
- Entities: User, Wallet, Transaction, WalletMetric, Watchlist
- Alchemy client: EVM (incoming + outgoing transfers) & Solana
- DexScreener client (liquidity, mcap, pair age, price change)
- Moralis client (EVM historical price by block + Solana current price)
- Birdeye client (Solana historical price by unix timestamp)
- Gemini AI client (natural language wallet insights)
- Jupiter client (Solana swap quotes)
- 1inch client (EVM swap quotes)
- Multi-chain support: Ethereum, Arbitrum, Optimism, Base, Polygon, Solana
- Chain registry (constants/chains.go)
- Real scoring: win rate, profit consistency, entry timing, token quality, trade discipline
- Risk exposure (informational only, not in score formula)
- FIFO buy-sell matching for PnL (realized + unrealized)
- 3-tier recommendations: full_follow, conditional_follow, avoid
- Conditional follow with AI-generated conditions (liquidity, pair age, mcap, volume)
- Traded tokens in response (contract address, symbol, PnL, status)
- AI insight (Gemini, beginner-friendly)
- Firebase Auth (Google login, backend token verification middleware)
- User system (entity, repository, service, find-or-create from Firebase)
- Watchlist (follow/unfollow wallet, conditions, email preference)
- DI container + router pattern (public + protected routes)
- Interfaces with I prefix
- Error handling (AppError + output envelope)
- Validation (go-playground/validator + ParseAndValidateBody)
- Dockerfile, Docker Compose, Makefile, .gitignore
- Migrations (auto-migrate, reset, seed)
- Documentation (comments on all files)

## 🔲 Todo — Hackathon Priority

- [ ] WebSocket hub (manage connections per user, push notifications)
- [ ] Wallet monitor (background job: poll blockchain, detect new trades, check conditions, notify)
- [ ] Email notification service (SendGrid/Resend)
- [ ] Smart contracts — Fee Router (swap fee collection) + On-chain Wallet Score (verifiable scoring)
- [ ] Frontend (Next.js: wallet input, scoring dashboard, swap UI, follow button, notifications)
- [ ] Wallet connect (Phantom + MetaMask)
- [ ] Tests (minimal)

## 📋 Post-Hackathon

- Auto copy-trade (execute swap on behalf of user)
- Smart contract fee router
- Rate limiting
- Caching (Redis)
- Pagination for transaction history
- Logging & monitoring
