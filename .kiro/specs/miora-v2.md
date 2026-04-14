# Spec: Miora AI V2 — Trading Reputation Protocol + AI Agent for Base

## Requirements

### R1: Connect Frontend to Backend API
Replace all dummy data in frontend with real API calls to backend endpoints.
- Analyze page: call POST /api/wallets/analyze, display real scoring results
- Watchlist page: call GET /api/watchlist, display real followed wallets
- Watchlist detail: call GET /api/wallets/:address, display stored analysis
- Notifications: connect WebSocket for real-time trade alerts
- Auth: integrate Firebase Auth (Google login) with GET /api/auth/me

### R2: Cleanup Solana/V1 Code
Remove all Solana-specific code per #[[file:CLEANUP.md]]:
- Delete backend clients: solana.go, birdeye.go, jupiter.go
- Delete interfaces/birdeye.go
- Remove Solana references from config, constants, services, router
- Remove Solana references from frontend (chain selectors, tokens, dummy data, assets)
- Delete contracts/svm/ directory
- Delete Solana Bruno API docs

### R3: Deploy EAS Attestation on Base Sepolia
Publish wallet trading reputation scores on-chain via Ethereum Attestation Service:
- Create EAS schema: score, recommendation, metrics, timestamp, wallet address
- After wallet analysis, publish attestation to Base Sepolia
- Add attestation UID to WalletMetric entity
- API endpoint: GET /api/reputation/:address — return attestation data
- Verify attestation readable on BaseScan Sepolia

### R4: Basic AgentKit Proof of Concept
Integrate Coinbase AgentKit for autonomous trading:
- Create Agentic Wallet on Base Sepolia
- Agent config: budget, max_per_trade, risk_tolerance, conditions
- Agent monitors top-scored wallets → evaluates trade → checks conditions
- If all pass → execute swap via Agentic Wallet on Base Sepolia
- Notify user of agent action via WebSocket

### R5: x402 Reputation API
Monetize reputation queries via x402 micropayments:
- Endpoint: GET /api/reputation/query — returns score
- Requires x402 USDC micropayment on Base
- Test with x402 client on Base Sepolia

### R6: Update Frontend for V2 Narrative
- Update landing page copy: "Trading Reputation Protocol" not "Wallet Analyzer"
- Update features section to reflect V2 (EAS, AgentKit, x402)
- Add Agent setup page (/agent)
- Add Agent dashboard (status, trade history, pause/resume)
- Update chains section: Base-first, remove Solana

### R7: Base Batches Application
- Complete all [FILL] fields in #[[file:BASE_BATCHES_APPLICATION.md]]
- Record 1-minute founder video
- Make GitHub repo public
- Submit before April 27, 2026

---

## Design

### Architecture (V2)
```
Frontend (Next.js) → Backend (Go + Fiber) → External APIs (Alchemy, DexScreener, Moralis, Gemini)
                                           → On-chain (EAS Attestation, Agentic Wallet, x402)
                                           → Database (PostgreSQL)
                                           → WebSocket (real-time notifications)
```

### New Backend Components
| Component | File | Purpose |
|---|---|---|
| EAS Client | clients/eas.go | Publish attestations to Base Sepolia via EAS SDK |
| AgentKit Client | clients/agentkit.go | Interact with Coinbase AgentKit for autonomous trading |
| x402 Middleware | middleware/x402.go | Verify x402 micropayment before serving reputation data |
| Reputation Handler | handlers/reputation.go | GET /reputation/:address, GET /reputation/query |
| Agent Handler | handlers/agent.go | POST /agent/start, PUT /agent/config, POST /agent/pause, GET /agent/status |
| Agent Service | services/agent.go | Monitor top wallets → evaluate → execute via AgentKit |
| Agent Entity | entities/agent_config.go | budget, max_per_trade, risk_tolerance, conditions, status |

### New Frontend Pages
| Page | Route | Purpose |
|---|---|---|
| Agent Setup | /agent | Configure AI trading agent (budget, rules, conditions) |
| Agent Dashboard | /agent/dashboard | Monitor agent actions, trade history, pause/resume |

### Database Changes
| Entity | Change |
|---|---|
| WalletMetric | Add: attestation_uid (string, nullable) |
| AgentConfig (NEW) | user_id, budget, max_per_trade, risk_tolerance, conditions (JSON), status (active/paused), created_at |
| AgentTrade (NEW) | agent_config_id, wallet_address, token_address, token_symbol, amount, tx_hash, risk_assessment, created_at |

---

## Tasks

### Phase 1: Cleanup & Connect (Priority #1)
- [ ] Task 1.1: Execute CLEANUP.md — remove all Solana code from backend
- [ ] Task 1.2: Execute CLEANUP.md — remove all Solana code from frontend
- [ ] Task 1.3: Delete contracts/svm/ and Solana Bruno API docs
- [ ] Task 1.4: Verify backend compiles: `cd backend && go build ./...`
- [ ] Task 1.5: Verify frontend builds: `cd frontend && npm run build`
- [ ] Task 1.6: Create API client in frontend (lib/api.ts)
- [ ] Task 1.7: Connect analyze page to POST /api/wallets/analyze
- [ ] Task 1.8: Connect watchlist page to GET /api/watchlist
- [ ] Task 1.9: Connect watchlist detail to GET /api/wallets/:address
- [ ] Task 1.10: Connect WebSocket for real-time notifications
- [ ] Task 1.11: Integrate Firebase Auth in frontend (login flow)

### Phase 2: EAS Attestation (Priority #2)
- [ ] Task 2.1: Research EAS SDK integration in Go
- [ ] Task 2.2: Create EAS schema on Base Sepolia
- [ ] Task 2.3: Implement clients/eas.go
- [ ] Task 2.4: Add attestation_uid to WalletMetric entity + migration
- [ ] Task 2.5: Publish attestation after wallet analysis in services/wallet.go
- [ ] Task 2.6: Implement handlers/reputation.go (GET /reputation/:address)
- [ ] Task 2.7: Register routes in http/reputation.go
- [ ] Task 2.8: Verify attestation on BaseScan Sepolia

### Phase 3: AgentKit (Priority #3)
- [ ] Task 3.1: Research Coinbase AgentKit Go/TypeScript SDK
- [ ] Task 3.2: Create Agentic Wallet on Base Sepolia
- [ ] Task 3.3: Implement entities/agent_config.go + entities/agent_trade.go
- [ ] Task 3.4: Implement services/agent.go (monitor → evaluate → execute)
- [ ] Task 3.5: Implement handlers/agent.go (start, config, pause, status)
- [ ] Task 3.6: Register routes in http/agent.go
- [ ] Task 3.7: Agent setup page in frontend (/agent)
- [ ] Task 3.8: Agent dashboard in frontend (/agent/dashboard)

### Phase 4: Frontend V2 Updates (Priority #4)
- [ ] Task 4.1: Update landing page copy for V2 narrative
- [ ] Task 4.2: Update features section (EAS, AgentKit, x402)
- [ ] Task 4.3: Update chains section (Base-first, remove Solana)
- [ ] Task 4.4: Add "Agent" nav item
- [ ] Task 4.5: Update analyze-form.tsx (remove Solana chain option)

### Phase 5: Submission (Priority #5)
- [ ] Task 5.1: Complete BASE_BATCHES_APPLICATION.md [FILL] fields
- [ ] Task 5.2: Record 1-minute founder video
- [ ] Task 5.3: Deploy frontend to Vercel
- [ ] Task 5.4: Deploy backend to Railway/Render
- [ ] Task 5.5: Make GitHub repo public
- [ ] Task 5.6: Submit to Base Batches before April 27
