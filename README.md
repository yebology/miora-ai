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
- Scores published on-chain via Ethereum Attestation Service (EAS) on Base
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
- Powered by Coinbase AgentKit + Agentic Wallets on Base
- Pause, adjust, or stop anytime

### 💰 x402 Reputation API
- Other protocols and AI agents can query Miora's reputation scores
- Pay-per-request via x402 micropayments (USDC on Base)
- No API keys needed — just connect wallet and pay

### 🔐 Authentication
- Google login via Firebase Auth
- Wallet connect (MetaMask for EVM via wagmi/viem)

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
│  │              On-chain Layer (Base)                       │  │
│  │  EAS Attestation · Agentic Wallet · x402 Payments       │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              External APIs                               │  │
│  │  Alchemy · DexScreener · Moralis · Gemini AI             │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Database (PostgreSQL)                        │  │
│  │  Users · Wallets · Transactions · Metrics ·              │  │
│  │  Watchlist · Notifications · Agent Configs               │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 🌐 Supported Chains

| Chain | Wallet Analysis | Swap Quotes | Historical Price |
|-------|:-:|:-:|:-:|
| **Base** | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Ethereum | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Arbitrum | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Optimism | ✅ | ✅ (1inch) | ✅ (Moralis) |
| Polygon | ✅ | ✅ (1inch) | ✅ (Moralis) |

Base is the primary chain. Other EVM chains are supported for wallet analysis.

---

## ⚙️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 16, Tailwind CSS v4, shadcn/ui, TypeScript, wagmi, viem |
| Backend | Go, Fiber, GORM, WebSocket |
| Database | PostgreSQL |
| Auth | Firebase Auth (Google) |
| AI | Google Gemini (gemini-2.0-flash) |
| Blockchain Data | Alchemy, DexScreener, Moralis |
| On-chain | EAS (Ethereum Attestation Service) on Base |
| Agent | Coinbase AgentKit + Agentic Wallets |
| Payments | x402 protocol (micropayments) |
| DEX Aggregation | 1inch (EVM) |
| Infra | Docker, Docker Compose |
| Email | Resend |
| API Testing | Bruno |

---

## 🧩 Project Structure

```
├── backend/
│   ├── app/
│   │   ├── clients/        # External API clients (Alchemy, DexScreener, Moralis, Gemini, 1inch)
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
│   ├── migrations/         # Database migrations
│   ├── router/             # DI container + route setup
│   ├── utils/              # Shared utilities
│   ├── pkg/                # Shared packages (AppError)
│   └── main.go             # Entry point
├── frontend/
│   ├── app/                # Next.js App Router pages
│   │   ├── page.tsx        # Landing page
│   │   ├── analyze/        # Wallet analysis page
│   │   ├── watchlist/      # Watchlist + detail pages
│   │   ├── swap/           # Swap page (placeholder)
│   │   └── login/          # Login page (placeholder)
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
│   └── lib/                # Utilities
├── contracts/
│   └── evm/                # EVM smart contracts (Foundry)
├── Makefile                # Dev commands
├── PROGRESS.md             # Development progress tracker
├── MIORA_V2_CONCEPT.md     # V2 concept document
└── README.md
```

---

## 🧭 How to Run

### 📦 Prerequisites
- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- Alchemy, Moralis, Gemini, 1inch API keys
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
make run-be
```

### 🌐 5. Run Frontend

```bash
make run-fe
```

---

## 📡 API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/health` | Health check |
| POST | `/api/wallets/analyze` | Analyze a wallet address |
| POST | `/api/wallets/regenerate-insight` | Regenerate AI insight with different tone |
| GET | `/api/wallets/:address` | Get stored analysis |
| POST | `/api/swap/quote` | Get swap quote (1inch) |

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
| `ws://host/ws?user_id=ID` | Real-time trade notifications |

### V2 (Planned)
| Method | Endpoint | Description |
|--------|----------|------------|
| GET | `/api/reputation/:address` | Get on-chain reputation attestation |
| GET | `/api/reputation/query` | Query reputation via x402 payment |
| POST | `/api/agent/start` | Start AI trading agent |
| PUT | `/api/agent/config` | Update agent configuration |
| POST | `/api/agent/pause` | Pause/resume agent |
| GET | `/api/agent/status` | Get agent status + trade history |

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
