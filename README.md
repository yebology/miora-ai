# 🧠 Miora AI

> **"See beyond the wallet."**

**Miora AI** is an AI-powered DEX aggregator with wallet intelligence — helping users understand blockchain wallet activity, get actionable trading recommendations, and swap tokens across multiple chains, all in one platform.

Instead of presenting raw on-chain data, Miora AI transforms complex blockchain activity into human-readable insights and decision support powered by AI.

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

Example:
> "This wallet is a disciplined trader that focuses on 3-4 tokens with consistent 20-30% gains. Safe to follow."

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

### 🔔 Smart Alerts & Watchlist
- Follow wallets and get real-time notifications when they trade
- Set custom conditions: "Only notify me if token liquidity > $100k and pair age > 6 hours"
- Email notifications with AI-generated token insights

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
6. 🔔 **Notify** — When the followed wallet trades, user gets notified (WebSocket + email) with token analysis
7. 🔄 **Trade** — User can swap tokens directly from Miora via DEX aggregator

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
| Frontend | Next.js, TailwindCSS, TypeScript |
| Backend | Go, Fiber, GORM |
| Database | PostgreSQL |
| Auth | Firebase Auth (Google) |
| AI | Google Gemini (gemini-2.0-flash) |
| Blockchain Data | Alchemy, DexScreener, Moralis, Birdeye |
| DEX Aggregation | Jupiter (Solana), 1inch (EVM) |
| Smart Contracts | Anchor (Solana), Foundry (EVM) |
| Infra | Docker, Docker Compose |

---

## 🧩 Project Structure

```
├── backend/
│   ├── app/
│   │   ├── clients/        # External API clients (Alchemy, DexScreener, Moralis, Birdeye, Gemini, Jupiter, 1inch)
│   │   ├── dto/            # Data transfer objects (requests, responses, prompts)
│   │   ├── entities/       # Database models (User, Wallet, Transaction, WalletMetric, Watchlist)
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── http/           # Route registration per domain
│   │   ├── interfaces/     # Service & repository contracts
│   │   ├── middleware/      # Firebase auth middleware
│   │   ├── output/         # Standardized API response (success/error envelope)
│   │   ├── repositories/   # Database access layer
│   │   ├── services/       # Business logic (wallet analysis, scoring, AI, swap, watchlist)
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
├── contracts/
│   ├── svm/                # Solana smart contracts (Anchor)
│   └── evm/                # EVM smart contracts (Foundry)
├── frontend/               # Next.js frontend
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
git clone https://github.com/your-username/miora-ai.git
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

---

## 📡 API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/health` | Health check |
| POST | `/api/wallets/analyze` | Analyze a wallet address |
| GET | `/api/wallets/:address` | Get stored analysis |
| POST | `/api/swap/quote` | Get swap quote (Jupiter/1inch) |

### Protected (Firebase Auth)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/auth/me` | Get/create current user |
| POST | `/api/watchlist/follow` | Follow a wallet |
| DELETE | `/api/watchlist/:address` | Unfollow a wallet |
| GET | `/api/watchlist` | List followed wallets |

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
