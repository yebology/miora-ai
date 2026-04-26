# 🧠 Miora AI

> **Score any trader on Base. Let AI ride the winners for you.**

🎬 **[Watch Founder Video](https://drive.google.com/file/d/1NrBd72H5XUzvQeBCYdTPXT7oHo86flpK/view?usp=sharing)** | 🖥️ **[Watch Demo Video](https://drive.google.com/file/d/1V73h7VNh0EZggx8mVDi0BXFlDTV8oL_1/view?usp=sharing)**

**Miora AI** is a trading reputation protocol on Base that analyzes any wallet's trading behavior, publishes scores on-chain as EAS attestations, and runs AI bots that trade autonomously based on the best wallets — with your rules, your budget, your conditions.

Instead of raw on-chain data, Miora transforms complex blockchain activity into actionable intelligence: a score, a recommendation, and a bot that acts on it.

---

## ✨ Overview

Every day, thousands of wallets trade on Base. Some are great traders. Most aren't. But there's no way to tell the difference — until now.

Miora combines three layers into one cohesive product:

1. 🧠 **Trading Reputation Protocol** — Analyze any wallet, compute a multi-factor score, publish it on-chain via EAS attestation. Other protocols can read and use this score.
2. 🔔 **Smart Follow + AI Alerts** — Follow top-scored wallets, get real-time notifications with AI risk assessment when they trade.
3. 🤖 **AI Trading Bots** — Per-wallet bots that monitor a target wallet's buys and sells, evaluate trades, and execute swaps on your behalf via Coinbase AgentKit — with your budget and conditions.

---

## 🎯 Key Features

### 🏆 On-chain Trading Reputation (EAS)
- Multi-factor scoring engine: win rate, profit consistency, entry timing, token quality, trade discipline
- FIFO buy-sell matching for accurate PnL calculation
- 3-tier recommendation: Full Follow (80-100), Conditional Follow (40-79), Avoid (<40)
- Scores published on-chain via Ethereum Attestation Service (EAS) on Base Sepolia
- Queryable by any protocol, agent, or dApp

### 🤖 AI-Powered Insights
- Google Gemini translates scoring data into beginner-friendly explanations
- AI risk assessment per trade notification — evaluates token liquidity, market cap, pair age before alerting
- Supports multiple tones: simple, eli5, custom prompt

### 🎯 Smart Recommendations & Conditional Follow
- Dynamic condition thresholds computed from wallet's own trading data
- Conditions: minimum liquidity, pair age, market cap, 24h volume
- Users choose which conditions to activate

### 🔔 Real-time Smart Alerts
- WebSocket notifications when followed wallets trade
- Each alert includes AI risk assessment
- Notification history saved to database

### 🤖 AI Trading Bots (AgentKit)
Two bot types, each targeting specific trading strategies:

**Wallet Bot** — Copy one specific wallet's trades (buys AND sells)
- User selects a wallet from their watchlist → conditions auto-filled from analyze result
- Bot monitors that wallet and mirrors its buys and sells
- User sets: budget, max per trade, min score, conditions
- Powered by Coinbase AgentKit + Agentic Wallets on Base Sepolia

**Consensus Bot** — Trade when multiple high-score wallets agree
- Scans all Miora-analyzed wallets on Base
- Trades when multiple wallets buy the same token within a time window
- User sets: budget, max per trade, min score, consensus threshold, time window
- Higher confidence trades based on crowd intelligence

Both bot types:
- Evaluate every trade through conditions from analyze result + AI risk assessment before executing
- Can be paused, adjusted, or stopped anytime
- Track all trades (executed/skipped/failed) with reasons

### 🔐 Authentication
- Wallet-based auth via MetaMask connect (wagmi/viem)
- `X-Wallet-Address` header for API authentication
- No Firebase — wallet connect only

---

## 🧩 System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js 16)                      │
│  Analyze → Dashboard → Bot Setup → Notifications             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  Backend (Go + Fiber)                         │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │ Scoring      │  │ Smart Follow │  │ AI Trading Bots   │  │
│  │ Engine       │  │ + Alerts     │  │ (AgentKit)        │  │
│  └──────┬───────┘  └──────────────┘  └─────────┬─────────┘  │
│         │                                       │            │
│  ┌──────▼───────────────────────────────────────▼─────────┐  │
│  │              On-chain Layer (Base Sepolia)               │  │
│  │  EAS Attestation · CDP Server Wallet · MockUSDT         │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              External APIs                               │  │
│  │  Alchemy · DexScreener · Moralis · Gemini AI             │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Agent Sidecar (Python + AgentKit)            │  │
│  │  Coinbase AgentKit · CDP Server Wallet · Swap Execution  │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Database (PostgreSQL)                        │  │
│  │  Users · Wallets · Transactions · Metrics ·              │  │
│  │  Watchlist · Notifications · Bot Configs · Bot Trades    │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## ⚙️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 16, Tailwind CSS v4, shadcn/ui, TypeScript, Zod, wagmi, viem |
| Backend | Go, Fiber, GORM, WebSocket |
| Agent Sidecar | Python, FastAPI, Coinbase AgentKit |
| Database | PostgreSQL |
| Auth | Wallet-based (MetaMask connect, X-Wallet-Address header) |
| AI | Google Gemini (gemini-2.5-flash) |
| Blockchain Data | Alchemy, DexScreener, Moralis |
| On-chain | EAS (Ethereum Attestation Service) on Base Sepolia |
| Agent | Coinbase AgentKit + CDP Server Wallet |
| Token | MockUSDT (ERC-20, 6 decimals) on Base Sepolia |
| Infra | Docker, Docker Compose |

---

## 🧩 Project Structure

```
├── backend/
│   ├── app/
│   │   ├── clients/        # External API clients (Alchemy, DexScreener, Moralis, Gemini, EAS, AgentKit)
│   │   ├── dto/            # Data transfer objects (requests, responses, prompts)
│   │   ├── entities/       # Database models (User, Wallet, Transaction, WalletMetric, Watchlist, Notification, AgentConfig, AgentTrade)
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── http/           # Route registration per domain
│   │   ├── interfaces/     # Service & repository contracts
│   │   ├── middleware/      # Wallet auth middleware (X-Wallet-Address)
│   │   ├── output/         # Standardized API response (success/error envelope)
│   │   ├── repositories/   # Database access layer
│   │   ├── services/       # Business logic (wallet, scoring, AI, watchlist, monitor, agent)
│   │   └── ws/             # WebSocket hub
│   ├── cmd/                # CLI commands (seed, reset, register-schema)
│   ├── config/             # Environment config loader
│   ├── constants/          # Constants (chains, errors, success messages, agent defaults)
│   ├── migrations/         # Database migrations
│   ├── router/             # DI container + route setup
│   ├── utils/              # Shared utilities
│   ├── pkg/                # Shared packages (AppError)
│   └── main.go             # Entry point
├── agent/
│   ├── main.py             # AgentKit sidecar (FastAPI + Coinbase AgentKit)
│   └── requirements.txt    # Python dependencies
├── frontend/
│   ├── api/                # API layer (per-module, Zod-validated)
│   │   ├── client.ts       # Core fetch wrapper + response envelope validation
│   │   ├── wallet/         # Wallet analysis API (validation.ts + connector.ts)
│   │   ├── watchlist/      # Watchlist CRUD API (validation.ts + connector.ts)
│   │   ├── reputation/     # Reputation API (validation.ts + connector.ts)
│   │   └── agent/          # Bot management API (validation.ts + connector.ts)
│   ├── app/                # Next.js App Router pages
│   │   ├── page.tsx        # Landing page
│   │   ├── analyze/        # Wallet analysis page
│   │   ├── watchlist/      # Watchlist + detail pages
│   │   ├── agent/          # Bot management pages
│   │   └── login/          # Login page
│   ├── components/
│   │   ├── ui/             # shadcn/ui components
│   │   ├── layout/         # Navbar, Footer, ThemeToggle
│   │   ├── landing/        # Landing page sections
│   │   ├── analyze/        # Analyze page components
│   │   ├── watchlist/      # Watchlist components
│   │   ├── agent/          # Bot config + status components
│   │   └── providers/      # Theme, Auth, Web3 providers
│   ├── constants/          # Static data + dummy data
│   ├── hooks/              # Custom hooks
│   ├── types/              # TypeScript types
│   └── lib/                # Utilities (cn helper)
├── contracts/
│   ├── src/MockUSDT.sol     # MockUSDT ERC-20 token (6 decimals)
│   └── script/              # Foundry deploy scripts
├── Makefile                # Dev commands
├── README.md               # Project overview
├── USER_STORIES.md          # User stories with scenarios
├── USER_NONTECHNICAL_FLOW.md # Non-technical user flow
└── USER_TECHNICAL_FLOW.md   # Technical architecture flow
```

---

## 🧭 How to Run

### 📦 Prerequisites
- Docker & Docker Compose
- MetaMask wallet (for authentication, connected to Base Sepolia)
- API keys: Alchemy, Moralis, Gemini
- Coinbase Developer Platform (CDP) credentials (API key + Wallet Secret)

### 🔨 1. Clone Repository

```bash
git clone https://github.com/yebology/miora-ai.git
cd miora-ai
```

### 🔐 2. Configure Environment

Fill in the `.env` files:
- `backend/.env` — DB, Alchemy, Moralis, Gemini, EAS, MockUSDT
- `agent/.env` — CDP API key, CDP wallet secret
- `frontend/.env` — API URL, Reown project ID, MockUSDT address

### 🚀 3. Run Everything (Docker)

```bash
make docker-up
```

This builds and starts all 4 services:
- **db** — PostgreSQL on port 5432
- **backend** — Go API on port 8082
- **agent** — AgentKit sidecar on port 8090
- **frontend** — Next.js on port 3002

### 📋 4. Register EAS Schema (one-time)

```bash
make register-schema
# Copy the printed schema UID to EAS_SCHEMA_UID in backend/.env
```

### 🌱 5. Seed Demo Data (optional)

```bash
make db-seed
```

### 🌐 6. Open App

Visit `http://localhost:3002` and connect MetaMask (Base Sepolia network).

---

## 📡 API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/health` | Health check |
| POST | `/api/wallets/analyze` | Analyze a wallet address on Base |
| POST | `/api/wallets/regenerate-insight` | Regenerate AI insight with different tone |
| GET | `/api/wallets/:address` | Get stored analysis |
| GET | `/api/reputation/:address` | Get on-chain reputation attestation |

### Protected (Wallet Auth — X-Wallet-Address header)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/auth/me` | Get/create current user |
| POST | `/api/watchlist/follow` | Follow a wallet with conditions |
| PUT | `/api/watchlist/:address` | Update conditions / notification preference |
| DELETE | `/api/watchlist/:address` | Unfollow a wallet |
| GET | `/api/watchlist` | List followed wallets |
| POST | `/api/agent/bots` | Create a new bot (wallet or consensus type) |
| GET | `/api/agent/bots` | List all user's bots |
| GET | `/api/agent/bots/:id` | Get bot details + status |
| PUT | `/api/agent/bots/:id` | Update bot configuration |
| POST | `/api/agent/bots/:id/start` | Start a bot |
| POST | `/api/agent/bots/:id/pause` | Pause a bot |
| DELETE | `/api/agent/bots/:id` | Delete a bot |
| GET | `/api/agent/bots/:id/trades` | Get bot's trade history |

### WebSocket
| Endpoint | Description |
|----------|------------|
| `ws://host/ws?wallet_address=ADDR` | Real-time trade notifications |

---

## 🗺️ Roadmap

### ✅ Hackathon (Done)
- EAS schema registered + attestation working on Base Sepolia
- Wallet bot + consensus bot working via Coinbase AgentKit
- Frontend connected to backend API (Zod-validated)
- MockUSDT deployed on Base Sepolia for bot deposits
- CDP Server Wallet integration for agent trading

### Post-Hackathon
- Deploy to Base mainnet
- Consensus bot as premium feature (revenue stream)
- Full DEX integration for bot swaps (Aerodrome/Uniswap)
- Reputation leaderboard
- Multi-bot portfolio tracking

### Scale
- Multi-chain expansion (Ethereum, Arbitrum, Optimism)
- Reputation score marketplace (protocols subscribe)
- Advanced bot strategies (cross-wallet pattern detection, sentiment analysis)
- Mobile app

---

## 🔥 Why Miora?

| Existing Tools | Miora AI |
|------|--------|
| Show data | Show decisions |
| Charts & numbers | "Follow this wallet" or "Avoid" |
| Analytics only | Analytics + autonomous trading bots |
| Off-chain scores | On-chain reputation via EAS |
| No composability | Other protocols can query Miora scores |
| For advanced traders | For everyone |

---

## 🤝 Contributors

🧑 **Yobel Nathaniel Filipus**
- 🐙 Github: [@yebology](https://github.com/yebology)
- 💼 LinkedIn: [View Profile](https://linkedin.com/in/yobelnathanielfilipus)

---

## ⚠️ Disclaimer

Miora AI provides informational insights only and does not constitute financial advice. AI bot trading involves risk. Users are responsible for their own trading decisions and bot configurations.

---

## 📄 License

MIT License
