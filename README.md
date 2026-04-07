# 🧠 Miora AI

> **"See beyond the wallet."**

**Miora AI** is an AI-powered DEX aggregator with wallet intelligence — helping users analyze blockchain wallets, get smart follow/avoid recommendations, receive AI-assessed trade alerts, and swap tokens across 6 chains through a single interface.

Instead of presenting raw on-chain data, Miora AI transforms complex blockchain activity into beginner-friendly insights and actionable decision support. Analyze any wallet, follow the good ones, get notified when they trade (with AI risk assessment), and act — all without leaving the platform.

---

## ✨ Overview

Blockchain data is transparent but overwhelming. Existing tools are built for advanced users — charts, numbers, and jargon that beginners can't understand. Miora AI solves this by combining:

- 🧠 **AI-Powered Wallet Intelligence** — Analyze any wallet, get a score, and understand trading behavior in plain language
- � **DEX Aggregator** — Swap tokens across Jupiter (Solana) and 1inch (EVM) with best route discovery
- 📊 **Smart Recommendations** — Full Follow, Conditional Follow, or Avoid — with AI-generated conditions
- 🔔 **Smart Alerts** — Follow wallets and get notified when they trade, filtered by your conditions
- 🌐 **Multi-Chain** — Ethereum, Arbitrum, Optimism, Base, Polygon, and Solana

---

## � Features

### 🔍 Wallet Intelligence
Analyze any wallet address across EVM and Solana chains. Get a comprehensive scoring based on:
- **Win Rate** — Percentage of profitable trades (realized + unrealized PnL)
- **Profit Consistency** — How stable the profits are across trades (standard deviation)
- **Entry Timing** — How early the wallet enters new tokens after launch
- **Token Quality** — Average market cap of tokens traded (logarithmic scale)
- **Trade Discipline** — How focused the wallet is (unique tokens vs total transactions)
- **Risk Exposure** — Percentage of low-liquidity tokens traded (informational)

### 🤖 AI-Powered Insights
The backend performs all the heavy analysis — fetching on-chain data, calculating PnL with FIFO buy-sell matching, and computing multi-factor scoring. The AI layer (Google Gemini) then takes these computed results and translates them into beginner-friendly, plain language explanations. AI does not analyze data itself — it narrates what the backend already calculated.

AI is also used to generate real-time risk assessments for trade notifications. When a followed wallet makes a trade, Gemini evaluates the token's market data (liquidity, market cap, pair age, price change) and provides a short risk opinion — helping users decide whether to act on the alert.

Example wallet insight:
> "This wallet is a disciplined trader that focuses on 3-4 tokens with consistent 20-30% gains. Safe to follow."

Example trade assessment:
> "⚠️ This token launched 45 minutes ago with only $8k liquidity. High risk entry — only follow if you're comfortable with potential loss."

### 🎯 Smart Recommendations
Three-tier recommendation system based on wallet score (0-100):

| Score | Recommendation | Action |
|-------|---------------|--------|
| 80-100 | ✅ Full Follow | Safe to follow — all trades shown with Buy button |
| 40-79 | ⚠️ Conditional Follow | Follow with conditions — AI suggests filters (liquidity, pair age, mcap) |
| < 40 | 🔴 Avoid | Do not follow — warning displayed |

### 🔄 DEX Aggregator
Swap tokens directly from Miora with best price routing:
- **Solana** → Jupiter (Raydium, Orca, Meteora, Lifinity, etc.)
- **EVM** → 1inch (Uniswap, SushiSwap, Curve, Balancer, etc.)

Users don't need to choose between Jupiter or 1inch — Miora automatically routes to the right aggregator based on the selected chain. One interface for all chains, no context switching between different DEX frontends.

### 🔔 Smart Alerts & Watchlist
- Follow wallets and get real-time notifications via WebSocket when they trade
- Each notification includes an AI risk assessment — Gemini evaluates the token and tells you if it's safe, risky, or dangerous
- Set custom conditions: "Only notify me if token liquidity > $100k and pair age > 6 hours"
- Notifications include trade details: token, amount, direction (buy/sell), timestamp, liquidity, market cap
- Notification history saved to database — never miss an alert even when offline
- Email notifications via Resend (async, non-blocking)

### 🔐 Authentication
- Google login via Firebase Auth
- Wallet connect (Phantom for Solana, MetaMask for EVM)

---

## 📋 How It Works

1. 🔍 **Analyze** — User inputs a wallet address and selects a chain
2. 📊 **Score** — Backend fetches on-chain data (Alchemy), enriches with market data (DexScreener), calculates PnL (Moralis/Birdeye), and generates a multi-factor score
3. 🤖 **Insight** — AI generates a beginner-friendly explanation of the wallet's trading behavior
4. 🎯 **Recommend** — System outputs Full Follow, Conditional Follow, or Avoid
5. 👀 **Follow** — User can follow the wallet with custom notification conditions
6. 🔔 **Notify** — When the followed wallet trades, user gets notified (in-app + email) with AI risk assessment
7. 🔄 **Trade** — User can swap tokens directly from Miora via the unified DEX aggregator

---

## 🧩 System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Frontend (Next.js)                  │
│  Wallet Input → Dashboard → Swap UI → Notifications     │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                    Backend (Go + Fiber)                   │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ Wallet   │  │ Swap     │  │ Auth     │  │Watchlist│ │
│  │ Analysis │  │ Quotes   │  │ Firebase │  │ + Alerts│ │
│  └────┬─────┘  └────┬─────┘  └──────────┘  └─────────┘ │
│       │              │                                   │
│  ┌────▼──────────────▼──────────────────────────────┐   │
│  │              External APIs                        │   │
│  │  Alchemy · DexScreener · Moralis · Birdeye       │   │
│  │  Jupiter · 1inch · Gemini AI                     │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Smart Contracts                      │   │
│  │  Fee Router (swap fees) · On-chain Score          │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Database (PostgreSQL)                 │   │
│  │  Users · Wallets · Transactions · Metrics ·       │   │
│  │  Watchlist                                        │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

---

## 🌐 Supported Chains

| Chain | Wallet Analysis | Swap Quotes | Historical Price |
|-------|:-:|:-:|:-:|
| Ethereum | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Arbitrum | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Optimism | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Base | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Polygon | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Solana | ✅ | ✅ (Jupiter) | ✅ (Birdeye) |

---

## ⚙️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 16, Tailwind CSS v4, shadcn/ui, TypeScript, next-themes |
| Backend | Go, Fiber, GORM, WebSocket |
| Database | PostgreSQL |
| Auth | Firebase Auth (Google) |
| AI | Google Gemini (gemini-2.0-flash) |
| Blockchain Data | Alchemy, DexScreener, Moralis, Birdeye |
| DEX Aggregation | Jupiter (Solana), 1inch (EVM) |
| Smart Contracts | Anchor (Solana), Foundry (EVM) |
| Infra | Docker, Docker Compose |
| Email | Resend |
| API Testing | Bruno |

---

## 🧩 Project Structure

```
├── backend/
│   ├── app/
│   │   ├── clients/        # External API clients (Alchemy, DexScreener, Moralis, Birdeye, Gemini, Jupiter, 1inch)
│   │   ├── dto/            # Data transfer objects (requests, responses, prompts)
│   │   ├── entities/       # Database models (User, Wallet, Transaction, WalletMetric, Watchlist, Notification)
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── http/           # Route registration per domain
│   │   ├── interfaces/     # Service & repository contracts
│   │   ├── middleware/      # Firebase auth middleware
│   │   ├── output/         # Standardized API response (success/error envelope)
│   │   ├── repositories/   # Database access layer
│   │   ├── services/       # Business logic (wallet analysis, scoring, AI, swap, watchlist, monitor)
│   │   └── ws/             # WebSocket hub
│   ├── cmd/                # CLI commands (seed, reset)
│   ├── config/             # Environment config loader
│   ├── constants/          # Constants (chains, errors, success messages)
│   ├── migrations/         # Database migrations (auto-migrate, reset, seed)
│   ├── router/             # DI container + route setup
│   ├── utils/              # Shared utilities (validator, math, helpers)
│   ├── pkg/                # Shared packages (AppError)
│   ├── main.go             # Entry point
│   ├── Dockerfile          # Multi-stage Docker build
│   └── docker-compose.yml  # PostgreSQL
├── frontend/
│   ├── app/                # Next.js App Router pages
│   │   ├── page.tsx        # Landing page
│   │   ├── analyze/        # Wallet analysis page
│   │   ├── watchlist/      # Watchlist + detail pages (/watchlist/[chain]/[address])
│   │   ├── swap/           # Swap page (placeholder)
│   │   └── login/          # Login page (placeholder)
│   ├── components/
│   │   ├── ui/             # shadcn/ui components (button, card, badge, dialog, etc.)
│   │   ├── layout/         # Navbar, Footer, ThemeToggle
│   │   ├── landing/        # Landing page sections (hero, features, how-it-works, chains, cta)
│   │   ├── analyze/        # Analyze page components (score ring, metric bars, conditions, AI insight, tokens table)
│   │   ├── watchlist/      # Watchlist components (wallet card, notification item)
│   │   └── providers/      # Theme provider
│   ├── constants/          # Static data (landing, nav, dummy data)
│   ├── hooks/              # Custom hooks (useAnimateOnScroll)
│   ├── types/              # TypeScript types (wallet, watchlist, api)
│   └── lib/                # Utilities (cn)
├── contracts/
│   ├── svm/                # Solana smart contracts (Anchor)
│   └── evm/                # EVM smart contracts (Foundry)
├── Makefile                # Dev commands
└── README.md
```

---

## 🧭 How to Run

### 📦 Prerequisites
- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- Alchemy, Moralis, Birdeye, Gemini, 1inch API keys
- Firebase project with Google sign-in enabled

### 🔨 1. Clone Repository

```bash
git clone https://github.com/yebology/miora-ai.git
cd miora-ai
```

### 🔐 2. Configure Environment

```bash
cp backend/.env.example backend/.env
# Fill in all API keys and Firebase credentials
```

### 🐘 3. Start Database

```bash
cd backend && docker compose up -d
```

### 🚀 4. Run Backend

```bash
cd backend && go run main.go
```

### 🌐 5. Run Frontend

```bash
cd frontend && npm install && npm run dev
```

---

## 🔑 Environment Variables

| Variable | Description |
|----------|------------|
| `APP_PORT` | Backend server port |
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_DB` | PostgreSQL database name |
| `DB_HOST` | Database host |
| `DB_PORT` | Database port |
| `ALCHEMY_API_KEY` | Alchemy API key (EVM + Solana RPC) |
| `MORALIS_API_KEY` | Moralis API key (EVM historical prices) |
| `BIRDEYE_API_KEY` | Birdeye API key (Solana historical prices) |
| `GEMINI_API_KEY` | Google Gemini API key (AI insights) |
| `ONEINCH_API_KEY` | 1inch API key (EVM swap quotes) |
| `FIREBASE_CREDENTIALS` | Path to Firebase service account JSON |
| `ALLOWED_ORIGINS` | CORS allowed origins |
| `SCORING_LIQUIDITY_THRESHOLD` | Min liquidity for risk exposure (USD) |
| `SCORING_ENTRY_TIMING_MAX_AGE` | Max pair age for entry timing (hours) |
| `SCORING_TOKEN_QUALITY_LOG_BASE` | Log base for token quality score |
| `RESEND_API_KEY` | Resend API key for email notifications (optional) |
| `RESEND_FROM_EMAIL` | Sender email for Resend (default: onboarding@resend.dev) |

---

## 📡 API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/health` | Health check |
| POST | `/api/wallets/analyze` | Analyze a wallet address |
| POST | `/api/wallets/regenerate-insight` | Regenerate AI insight with different tone |
| GET | `/api/wallets/:address` | Get stored analysis |
| POST | `/api/swap/quote` | Get swap quote (Jupiter/1inch) |

### Protected (Firebase Auth)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/auth/me` | Get/create current user |
| POST | `/api/watchlist/follow` | Follow a wallet with conditions |
| PUT | `/api/watchlist/:address` | Update conditions / notification preference |
| DELETE | `/api/watchlist/:address` | Unfollow a wallet |
| GET | `/api/watchlist` | List followed wallets |

### WebSocket
| Endpoint | Description |
|----------|------------|
| `ws://host/ws?user_id=ID` | Real-time trade notifications for followed wallets |

---

---

## 🧪 API Testing

API documentation and test collection is in `api-docs/` using [Bruno](https://www.usebruno.com/) format.

### Setup
1. Install Bruno: https://www.usebruno.com/
2. Open Bruno → Import Collection → select `api-docs/` folder
3. Update variables in `collection.bru` (baseUrl, walletAddress, firebaseToken)
4. Run requests

### Collection Structure
```
api-docs/
├── bruno.json              → Collection metadata
├── collection.bru          → Shared variables + auth config
├── health/                 → Health check
├── wallets/                → Analyze, Get Stored, Regenerate Insight
├── swap/                   → Quote (Solana + EVM)
├── auth/                   → Get current user (protected)
└── watchlist/              → Follow, Update, List, Unfollow (protected)
```

---

## 🔥 Key Differentiation

| Existing Tools | Miora AI |
|------|--------|
| Data-heavy dashboards | AI-powered natural language insights |
| Charts & numbers | Beginner-friendly explanations |
| Analytics only | Analytics + DEX trading in one platform |
| No recommendations | Smart Follow/Avoid recommendations with conditions |
| For advanced traders | For everyone |

---

## 🤝 Contributors

🧑 **Yobel Nathaniel Filipus**
- 🐙 Github: [@yebology](https://github.com/yebology)
- 💼 LinkedIn: [View Profile](https://linkedin.com/in/yobelnathanielfilipus)
- � Email: yobelnathaniel12@gmail.com

---

## ⚠️ Disclaimer

Miora AI provides informational insights only and does not constitute financial advice. Users are responsible for their own trading decisions.

---

## 📄 License

MIT License
