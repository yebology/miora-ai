# Miora AI — Technical Flow

## Architecture Overview

```
Frontend (Next.js 16)
    │
    ▼ HTTP + WebSocket
Backend (Go + Fiber)
    │
    ├── Alchemy (Base RPC) — fetch transactions
    ├── DexScreener — token pair data
    ├── Moralis — historical prices
    ├── Gemini AI — insights + risk assessment
    ├── EAS (Base Sepolia) — on-chain attestations
    ├── PostgreSQL — all data persistence
    └── WebSocket Hub — real-time notifications
    │
    ▼ HTTP (localhost:8090)
Agent Sidecar (Python + FastAPI + AgentKit)
    │
    └── Coinbase AgentKit → CDP Server Wallet → Base Sepolia
```

---

## Flow 1: Wallet Analysis

```
POST /api/wallets/analyze { address: "0x...", chain: "base" }
    │
    ▼ handlers/wallet.go → Analyze()
    │
    ▼ services/wallet.go → AnalyzeWallet()
    │
    ├── 1. clients/evm.go → Alchemy GetTransfers(address, limit, "base")
    │      Returns: []TransferData (hash, from, to, value, token, direction, block)
    │
    ├── 2. clients/dexscreener.go → GetTokenPairs(chain, contractAddress)
    │      Returns: []TokenPairData (liquidity, mcap, pairAge, volume, priceChange)
    │
    ├── 3. clients/moralis.go → GetHistoricalPrice(chain, address, block)
    │      Returns: TokenPrice (usdPrice at specific block)
    │
    ├── 4. services/wallet_helper.go → calculateTrades()
    │      FIFO buy-sell matching → []tradeResult (buyPrice, exitPrice, pnl%)
    │
    ├── 5. services/scoring.go → calculateMetrics()
    │      5 factors → WalletMetric (winRate, profitConsistency, entryTiming,
    │      tokenQuality, tradeDiscipline, finalScore, recommendation)
    │
    ├── 6. repositories/wallet.go → SaveMetric()
    │      Upsert to PostgreSQL
    │
    ├── 7. services/wallet_helper.go → buildConditions()
    │      Dynamic thresholds from wallet's own data (median liquidity, mcap, etc.)
    │
    ├── 8. services/ai.go → GenerateInsight()
    │      Gemini prompt with scoring data → plain language insight
    │
    └── 9. [ASYNC goroutine] clients/eas.go → Attest()
           ABI-encode (score, recommendation, totalTxns, chain)
           → Sign tx with attester private key
           → Send to EAS contract 0x4200...0021 on Base Sepolia
           → Parse Attested event log → extract attestation UID
           → repositories/wallet.go → SaveMetric() (update attestation_uid)
    │
    ▼ Response: WalletAnalysis JSON
```

---

## Flow 2: Follow Wallet

```
POST /api/watchlist/follow
Headers: X-Wallet-Address: 0xUserWallet
Body: { wallet_address, chain, recommendation, conditions, email_notify }
    │
    ▼ middleware/wallet_auth.go → extract wallet_address from header
    ▼ handlers/watchlist.go → Follow()
    ▼ services/user.go → FindOrCreateByWallet() → get/create user
    ▼ services/watchlist.go → Follow()
    ▼ repositories/watchlist.go → Create()
    │
    ▼ Saved to PostgreSQL: Watchlist { userID, walletAddress, chain, conditions }
```

---

## Flow 3: Real-time Notifications

```
services/monitor.go → Start() [goroutine, ticker 30s]
    │
    ├── monitor_helper.go → poll()
    │   └── getUniqueWatchedWallets() → all distinct wallets from watchlist table
    │
    ├── For each wallet:
    │   └── checkWallet(address, chain)
    │       ├── clients/evm.go → GetTransfers(address, 100, chain)
    │       ├── Compare tx count with lastTxCount map
    │       └── New txs? → notifyFollowers()
    │
    └── notifyFollowers(walletAddress, chain, tx)
        ├── repositories/watchlist.go → FindByWallet(address)
        ├── clients/dexscreener.go → GetTokenPairs() → token data
        ├── services/ai.go → GenerateTradeAssessment() → risk text
        ├── monitor_helper.go → meetsConditions() → check user's conditions
        │
        ├── ws/hub.go → SendToUser(userID, message) → WebSocket push
        └── repositories/notification.go → Create() → save to DB
```

---

## Flow 4: Wallet Bot (Auto-trade)

```
POST /api/agent/bots
Headers: X-Wallet-Address: 0xUserWallet
Body: { bot_type: "wallet", target_wallet_address, budget, max_per_trade, conditions }
    │
    ▼ handlers/agent.go → CreateBot()
    ▼ services/agent.go → CreateBot() → save AgentConfig to DB
    │
POST /api/agent/bots/:id/start → set status = "active"
    │
    ▼ services/agent_loop.go → Start() [goroutine, ticker 30s]
    │
    ├── poll() → FindActiveConfigs() → all bots with status "active"
    │
    ├── processConfig(config) → bot_type == "wallet"
    │   └── checkWalletForAgent(config, targetWallet)
    │       ├── clients/evm.go → GetTransfers(targetWallet, 10, chain)
    │       ├── Detect new txs (compare with lastTxCount)
    │       └── For each new tx → evaluateAndExecute()
    │
    └── evaluateAndExecute(config, wallet, tx)
        ├── Determine direction: tx.Direction "in" → buy, "out" → sell
        │
        ├── [BUY only] Check conditions:
        │   ├── clients/dexscreener.go → GetTokenPairs()
        │   ├── meetsAgentConditions() → liquidity, mcap, pairAge, volume
        │   └── Check budget: remaining >= maxPerTrade?
        │
        ├── services/ai.go → GenerateTradeAssessment() → risk text
        │
        ├── clients/agentkit.go → ExecuteSwap(token, symbol, amount, direction)
        │   └── HTTP POST http://localhost:8090/swap
        │       └── agent/main.py → AgentKit → Agentic Wallet → on-chain tx
        │
        ├── Update config: totalSpent += maxPerTrade, totalTrades++
        └── repositories/agent.go → CreateTrade() → save AgentTrade to DB
```

---

## Flow 5: Consensus Bot

```
POST /api/agent/bots
Body: { bot_type: "consensus", budget, max_per_trade, min_score, consensus_threshold, consensus_window_min }
    │
    ▼ Same create + start flow as wallet bot
    │
    ▼ services/agent_loop.go → processConfig(config) → bot_type == "consensus"
    │
    └── processConsensus(config)
        ├── repositories/wallet.go → FindAllWithMetrics(minScore)
        │   SQL: JOIN wallet_metrics WHERE final_score >= minScore
        │   Returns: all wallets in Miora DB with score >= threshold
        │
        ├── For each wallet:
        │   ├── clients/evm.go → GetTransfers()
        │   ├── Detect new buys
        │   └── Track: tokenBuyers[tokenAddress] = append(wallet)
        │
        ├── Check consensus:
        │   └── For each token where len(buyers) >= consensus_threshold:
        │       ├── Log: "CONSENSUS: 3 wallets bought TOKEN_X"
        │       ├── evaluateAndExecuteConsensus()
        │       │   ├── Check conditions (from config)
        │       │   ├── Check budget
        │       │   └── clients/agentkit.go → ExecuteSwap()
        │       └── Record trade with reason: "CONSENSUS: 3 wallets (avg score 85)"
```

---

## Flow 6: Query Reputation (On-chain)

```
GET /api/reputation/0x1234...
    │
    ▼ handlers/reputation.go → GetReputation()
    │
    ├── repositories/wallet.go → FindByAddress() → get wallet
    ├── repositories/wallet.go → GetMetric() → get score + attestation_uid
    │
    └── clients/eas.go → GetAttestation(uid)
        ├── ABI pack getAttestation(bytes32 uid)
        ├── ethclient.CallContract() → EAS contract 0x4200...0021
        ├── Unpack Attestation struct (uid, schema, time, attester, recipient, data)
        └── Unpack custom data: (uint8 score, string recommendation, uint32 totalTxns, string chain)
    │
    ▼ Response: { address, score, recommendation, attestation_uid, explorer_url, attester, timestamp }
```

---

## Flow 7: Authentication (Wallet-based)

```
Frontend: useAppKitAccount() → { address, isConnected }
    │
    ▼ When wallet connected:
    │   AuthProvider sets user = { walletAddress: address }
    │
    ▼ Protected API calls:
    │   Headers: X-Wallet-Address: 0xUserWallet
    │
    ▼ middleware/wallet_auth.go
    │   ├── Read X-Wallet-Address header
    │   ├── Validate: 42 chars, starts with 0x
    │   └── Set c.Locals("wallet_address", address)
    │
    ▼ handlers/*.go → getUserID()
        ├── Read c.Locals("wallet_address")
        ├── services/user.go → FindOrCreateByWallet()
        └── Return user.ID for DB operations
```

---

## Database Schema

```
users
  id, wallet_address (unique), email (optional), created_at, updated_at

wallets
  id, address (unique), chain, created_at, updated_at

transactions
  id, wallet_id (FK), hash, chain, from, to, value,
  token_symbol, contract_address, direction, block_number, timestamp
  (unique: hash + direction)

wallet_metrics
  id, wallet_id (unique FK), total_transactions,
  profit_consistency, win_rate, entry_timing,
  token_quality, trade_discipline, final_score, recommendation,
  ai_insight, ai_tone, ai_prompt,
  conditions_json, traded_tokens_json,
  attestation_uid, attestation_tx_hash, updated_at

watchlists
  id, user_id (FK), wallet_address, chain, recommendation,
  conditions (JSON), email_notify, created_at

notifications
  id, user_id (FK), wallet_address, chain, token_address, token_symbol,
  direction, value, liquidity, market_cap, ai_assessment, read, created_at

agent_configs
  id, user_id (FK), bot_type, target_wallet_address, target_wallet_chain,
  target_wallet_score, recommendation, budget, max_per_trade,
  conditions (JSON), status, agent_wallet_address,
  total_spent, total_trades, consensus_threshold, consensus_window_min,
  min_score, created_at, updated_at

agent_trades
  id, agent_config_id (FK), source_wallet, source_score,
  token_address, token_symbol, direction, amount_usd, tx_hash,
  status, reason, risk_assessment, created_at
```

---

## Key File Map

| Layer | Files |
|-------|-------|
| Config | `backend/config/config.go` |
| Constants | `backend/constants/chains.go`, `agent.go`, `error.go`, `limits.go` |
| Entities | `backend/app/entities/*.go` (8 models) |
| Interfaces | `backend/app/interfaces/*.go` (7 contracts) |
| Repositories | `backend/app/repositories/*.go` (4 repos) |
| Services | `backend/app/services/*.go` (scoring, wallet, ai, watchlist, monitor, agent, agent_loop) |
| Clients | `backend/app/clients/*.go` (evm, dexscreener, moralis, gemini, eas, agentkit) |
| Handlers | `backend/app/handlers/*.go` (auth, wallet, watchlist, reputation, agent) |
| Routes | `backend/app/http/*.go` + `backend/router/routes.go` |
| Middleware | `backend/app/middleware/wallet_auth.go` |
| Agent Sidecar | `agent/main.py` (FastAPI + Coinbase AgentKit + CDP Server Wallet) |
| Frontend Pages | `frontend/app/` (landing, analyze, watchlist, agent, login) |
| Frontend Components | `frontend/components/` (analyze, watchlist, agent, landing, layout, ui) |
