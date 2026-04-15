# 🧠 Miora AI

> **Score any trader on Base. Let AI ride the winners for you.**

**Miora AI** is a trading reputation protocol on Base that analyzes any wallet's trading behavior, publishes scores on-chain as EAS attestations, and runs an AI agent that trades autonomously based on the best wallets — with your rules, your budget, your risk limits.

Instead of raw on-chain data, Miora transforms complex blockchain activity into actionable intelligence: a score, a recommendation, and an agent that acts on it.

---

## ✨ Overview

Every day, thousands of wallets trade on Base. Some are great traders. Most aren't. But there's no way to tell the difference — until now.

Miora combines three layers into one cohesive product:

1. 🧠 **Trading Reputation Protocol** — Analyze any wallet, compute a multi-factor score, publish it on-chain via EAS attestation. Other protocols can read and use this score.
2. 🔔 **Smart Follow + AI Alerts** — Follow top-scored wallets, get real-time notifications with AI risk assessment when they trade.
3. 🤖 **AI Trading Agent** — An autonomous agent that monitors top wallets, evaluates trades, and executes swaps on your behalf via Coinbase AgentKit — with your budget and conditions.

---

## 🎯 Key Features

### 🏆 On-chain Trading Reputation (EAS)
- Multi-factor scoring engine: win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure
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
- Email notifications via Resend (async, non-blocking)
- Notification history saved to database

### 🤖 AI Trading Agent (AgentKit)
- Autonomous trading based on top-scored wallets
- User sets: budget, max per trade, risk tolerance, conditions
- Agent evaluates every trade through scoring engine + AI risk assessment before executing
- Powered by Coinbase AgentKit + Agentic Wallets on Base Sepolia
- Pause, adjust, or stop anytime

### 💰 x402 Reputation API
- Other protocols and AI agents can query Miora's reputation scores
- Pay-per-request via x402 micropayments (USDC on Base)
- No API keys needed — just connect wallet and pay

### 🔐 Authentication
- Google login via Firebase Auth
- Wallet connect (MetaMask via wagmi/viem)

---

## 🧩 System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js 16)                      │
│  Analyze → Dashboard → Agent Setup → Notifications           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  Backend (Go + Fiber)                         │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │ Scoring      │  │ Smart Follow │  │ AI Trading Agent  │  │
│  │ Engine       │  │ + Alerts     │  │ (AgentKit)        │  │
│  └──────┬───────┘  └──────────────┘  └─────────┬─────────┘  │
│         │                                       │            │
│  ┌──────▼───────────────────────────────────────▼─────────┐  │
│  │              On-chain Layer (Base Sepolia)               │  │
│  │  EAS Attestation · Agentic Wallet · x402 Payments       │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              External APIs                               │  │
│  │  Alchemy · DexScreener · Moralis · Gemini AI             │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Agent Sidecar (Python + AgentKit)            │  │
│  │  Coinbase AgentKit · Agentic Wallet · Swap Execution     │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Database (PostgreSQL)                        │  │
│  │  Users · Wallets · Transactions · Metrics ·              │  │
│  │  Watchlist · Notifications · Agent Configs · Agent Trades│  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## ⚙️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 16, Tailwind CSS v4, shadcn/ui, TypeScript, wagmi, viem |
| Backend | Go, Fiber, GORM, WebSocket |
| Agent Sidecar | Python, FastAPI, Coinbase AgentKit |
| Database | PostgreSQL |
| Auth | Firebase Auth (Google) |
| AI | Google Gemini (gemini-2.0-flash) |
| Blockchain Data | Alchemy, DexScreener, Moralis |
| On-chain | EAS (Ethereum Attestation Service) on Base Sepolia |
| Agent | Coinbase AgentKit + Agentic Wallets |
| Payments | x402 protocol (micropayments) |
| Infra | Docker, Docker Compose |
| Email | Resend |
| API Testing | Bruno |

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
│   │   ├── middleware/      # Firebase auth + x402 payment middleware
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
│   ├── app/                # Next.js App Router pages
│   │   ├── page.tsx        # Landing page
│   │   ├── analyze/        # Wallet analysis page
│   │   ├── watchlist/      # Watchlist + detail pages
│   │   └── login/          # Login page
│   ├── components/
│   │   ├── ui/             # shadcn/ui components
│   │   ├── layout/         # Navbar, Footer, ThemeToggle
│   │   ├── landing/        # Landing page sections
│   │   ├── analyze/        # Analyze page components
│   │   ├── watchlist/      # Watchlist components
│   │   └── providers/      # Theme, Auth, Web3 providers
│   ├── constants/          # Static data + dummy data
│   ├── hooks/              # Custom hooks
│   ├── types/              # TypeScript types
│   └── lib/                # Utilities (API client, helpers)
├── Makefile                # Dev commands
├── PROGRESS.md             # Development progress tracker
├── MIORA_V2_CONCEPT.md     # V2 concept document
└── README.md
```

---

## 🧭 How to Run

### 📦 Prerequisites
- Go 1.25+
- Python 3.10+ (for AgentKit sidecar)
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- Alchemy, Moralis, Gemini API keys
- Firebase project with Google sign-in enabled
- Coinbase Developer Platform (CDP) API keys (for AgentKit)

### 🔨 1. Clone Repository

```bash
git clone https://github.com/yebology/miora-ai.git
cd miora-ai
```

### 🔐 2. Configure Environment

```bash
cp backend/.env.example backend/.env
cp agent/.env.example agent/.env
# Fill in all API keys and credentials
```

### 🐘 3. Start Database

```bash
cd backend && docker compose up -d
```

### 📋 4. Register EAS Schema (one-time)

```bash
make register-schema
# Copy the printed schema UID to EAS_SCHEMA_UID in backend/.env
```

### 🚀 5. Run Backend

```bash
make run-be
```

### 🤖 6. Run Agent Sidecar

```bash
make setup-agent  # Install Python dependencies (first time only)
make run-agent    # Start AgentKit sidecar on port 8090
```

### 🌐 7. Run Frontend

```bash
make run-fe
```

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

### x402 Protected (USDC micropayment)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/reputation/query?address=0x...` | Query reputation score (requires x402 payment) |

### Protected (Firebase Auth)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/auth/me` | Get/create current user |
| POST | `/api/watchlist/follow` | Follow a wallet with conditions |
| PUT | `/api/watchlist/:address` | Update conditions / notification preference |
| DELETE | `/api/watchlist/:address` | Unfollow a wallet |
| GET | `/api/watchlist` | List followed wallets |
| GET | `/api/agent/status` | Get agent status + config |
| PUT | `/api/agent/config` | Update agent configuration |
| POST | `/api/agent/start` | Start AI trading agent |
| POST | `/api/agent/pause` | Pause agent |
| GET | `/api/agent/trades` | Get agent trade history |

### WebSocket
| Endpoint | Description |
|----------|------------|
| `ws://host/ws?user_id=ID` | Real-time trade notifications |

---

## 🔥 Why Miora?

| Existing Tools | Miora AI |
|------|--------|
| Show data | Show decisions |
| Charts & numbers | "Follow this wallet" or "Avoid" |
| Analytics only | Analytics + autonomous trading agent |
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

Miora AI provides informational insights only and does not constitute financial advice. AI agent trading involves risk. Users are responsible for their own trading decisions and agent configurations.

---

## 📄 License

MIT License
