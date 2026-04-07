# Miora AI — Progress

## ✅ Backend Done

- Clean architecture (handlers, services, repositories, interfaces, entities, dto)
- Config loader (env-based, no fallbacks, configurable scoring thresholds)
- Database (PostgreSQL via Docker Compose, GORM auto-migrate)
- Entities: User, Wallet, Transaction, WalletMetric, Watchlist, Notification
- Alchemy client: EVM (incoming + outgoing transfers) & Solana
- DexScreener client (liquidity, mcap, pair age, price change)
- Moralis client (EVM historical price by block + Solana current price)
- Birdeye client (Solana historical price by unix timestamp)
- Gemini AI client (natural language wallet insights)
- Jupiter client (Solana swap quotes)
- 1inch client (EVM swap quotes)
- Multi-chain support: Ethereum, Arbitrum, Optimism, Base, Polygon, Solana
- Chain registry (constants/chains.go)
- Configurable transaction fetch limits per chain type (constants/limits.go)
  - EVM: 10, 25, 50, 100 (default 10)
  - Solana: 20, 50, 100, 200 (default 20)
- Real scoring: win rate, profit consistency, entry timing, token quality, trade discipline
- Risk exposure (informational only, not in score formula)
- FIFO buy-sell matching for PnL (realized + unrealized)
- 3-tier recommendations: full_follow, conditional_follow, avoid
- Conditional follow with dynamic thresholds (computed from wallet's token data — median liquidity, mcap, volume, avg pair age)
- Traded tokens in response (contract address, symbol, PnL, status)
- AI insight with tone support (simple, eli5, custom prompt)
- Regenerate insight endpoint (POST /wallets/regenerate-insight)
- Firebase Auth (Google login, backend token verification middleware)
- User system (entity, repository, service, find-or-create from Firebase)
- Watchlist CRUD (follow, unfollow, update conditions, list)
- WebSocket hub (manage connections per user, push notifications)
- Wallet monitor (background polling, detect new trades, check conditions, notify via WebSocket + save to DB)
- Email notifications wired into monitor service (Resend API client + HTML email template + async dispatch)
- DI container + router pattern (public + protected routes)
- Interfaces with I prefix (including IEmailClient for Resend)
- Error handling (AppError + output envelope)
- Validation (go-playground/validator + ParseAndValidateBody)
- Dockerfile, Docker Compose, Makefile, .gitignore
- Migrations (auto-migrate, reset, seed)
- Documentation (comments on all files)
- Bruno API test collection (api-docs/) for all endpoints
- API Testing Flow guide (API_TESTING_FLOW.md)

## ✅ Frontend Done

- Next.js 16 + Tailwind v4 + shadcn/ui + TypeScript
- Dark/light mode (next-themes, default dark)
- Space Grotesk font
- Responsive navbar with mobile hamburger menu
- Landing page:
  - Hero section (gradient orbs background, shimmer text animation, CTA buttons)
  - Features section (6 cards with hover pop + glow effect)
  - How It Works (3 steps with icons, animated connector line, pulse rings)
  - Chains section (infinite marquee with logos, hover brand color glow)
  - CTA section
  - Footer
  - Scroll-triggered fade-in animations
- Analyze page:
  - Wallet address input + chain selector
  - Score ring with gradient color per tier
  - Recommendation badge (Full Follow / Conditional Follow / Avoid)
  - AI Insight card with regenerate (Simple / ELI5 / Custom prompt)
  - Scoring breakdown with animated bars + info tooltips (formula explanations)
  - Interactive conditions card (toggle checkboxes, descriptions, Follow CTA)
  - Trade summary stats (tokens traded, avg PnL, win/loss, realized)
  - Sortable traded tokens table with scrollable container
  - Error state + loading state
  - "Wallet already exists" confirmation modal
- Watchlist page:
  - Tabs: Wallets / Notifications (with unread badge)
  - Wallet cards (address, chain, recommendation, conditions, notify toggle, unfollow)
  - Notification items (buy/sell icon, token, amount, liquidity, mcap, time ago)
  - Empty states
- Watchlist detail page (/watchlist/[chain]/[address]):
  - Stored analysis view (reuses AnalysisResult component)
  - Activity tab (notification history per wallet)
  - Re-analyze in-place (confirmation modal + success modal)
- Placeholder pages: /swap, /login
- Types matching backend response shapes
- Dummy data for all pages

## 📋 Post-Hackathon

- Auto copy-trade (execute swap on behalf of user)
- Smart contract fee router
- Rate limiting
- Caching (Redis)
- Pagination for transaction history
- Logging & monitoring
- Auth guard modal (sign in required for follow, watchlist)
- Wallet connect guard modal (connect wallet required for swap)
- Dynamic transaction limit selector per chain (EVM vs Solana options)
- Follow CTA with auth check for full_follow (in analysis header)
- Info tooltips on Follow buttons explaining what follow does
- Dummy data for all 3 recommendation tiers (full_follow, conditional_follow, avoid)

## 🔲 Todo — Hackathon Priority

- [ ] Connect frontend to real backend API (replace dummy data)
- [ ] Smart contracts — Fee Router (swap fee collection) + On-chain Wallet Score
- [ ] Wallet connect (Phantom + MetaMask)
- [ ] Tests (minimal)
