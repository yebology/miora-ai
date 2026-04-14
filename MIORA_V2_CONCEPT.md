# Miora V2 — Trading Reputation Protocol + AI Agent for Base

## Konsep Baru dalam Satu Kalimat

**Miora adalah trading reputation protocol di Base yang menganalisis wallet, mempublikasikan score on-chain sebagai attestation, dan menjalankan AI agent yang trade otomatis berdasarkan pola wallet terbaik — dengan rules dan budget yang user tentukan.**

---

## Mengapa Harus Pivot dari V1?

### Miora V1 (Sekarang)
- Wallet analyzer + DEX aggregator + smart alerts
- User analyze → follow → dapat notifikasi → swap manual
- Semua off-chain, tidak ada on-chain state
- Multi-chain generic, tidak ada yang Base-specific
- Kategori: **analytics tool**

### Miora V2 (Baru)
- Trading reputation protocol + AI trading agent
- Score di-publish on-chain → protocol lain bisa baca → AI agent trade otomatis
- On-chain attestation via EAS di Base
- Deep Base integration (EAS + AgentKit + x402)
- Kategori: **infrastructure protocol + autonomous agent**

### Alasan Pivot

1. **Analytics tool sudah saturated** — Nansen, Arkham, DeBank, Cielo, GMGN sudah ada. Miora V1 adalah iterasi, bukan inovasi.
2. **Base 2026 strategy fokus ke 3 hal**: global markets, stablecoin payments, dan home for builders/agents. Miora V2 align dengan pilar 1 (markets — discover top traders) dan pilar 3 (agents — AI agent trading).
3. **BB003 cohort tidak punya trading intelligence layer** — Ada AI agent infra (Blockrun, Agently), lending (Credifi), privacy DEX (OPAL), tapi tidak ada yang bantu retail trader trade lebih smart.
4. **70% codebase V1 tetap terpakai** — Scoring engine, conditional follow, AI risk assessment, wallet monitoring semua jadi komponen V2. Bukan rewrite, tapi evolusi.

---

## Tiga Layer Produk

### Layer 1: Trading Reputation Protocol (On-chain)

**Apa ini:**
Miora menganalisis trading history wallet di Base, menghitung multi-factor score, lalu mempublikasikan score tersebut on-chain sebagai attestation menggunakan EAS (Ethereum Attestation Service).

**Bagaimana cara kerjanya:**
1. Miora fetch transaction history wallet dari Base via Alchemy
2. Scoring engine menghitung 6 faktor: win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure
3. Score (0-100) + metadata (tier, timestamp, jumlah trade yang dianalisis) di-publish ke Base sebagai EAS attestation
4. Attestation bisa dibaca oleh siapapun on-chain — protocol, dApp, AI agent

**Mengapa EAS:**
- EAS sudah live di Base mainnet dan Base Sepolia testnet
- Coinbase sendiri menggunakan EAS untuk Coinbase Verifications (KYC attestation on-chain)
- EAS adalah standard — bukan custom contract yang harus dipercaya orang
- Gratis untuk membuat attestation (hanya bayar gas, yang di Base sangat murah)
- Sudah ada SDK dan tooling yang mature

**Siapa yang pakai reputation score ini:**
| Consumer | Use Case |
|---|---|
| DeFi lending protocols | Assess borrower quality — "wallet ini trader bagus, kasih bunga lebih rendah" |
| Airdrop campaigns | Filter eligible wallets — "hanya wallet score > 70 yang dapat airdrop" |
| AI agents | Trust scoring — "sebelum copy trade wallet ini, cek Miora score dulu" |
| DEX platforms | Flag risky wallets — "wallet score < 30, tampilkan warning" |
| Users sendiri | Claim reputation — "saya trader dengan Miora score 90" |

**Monetization:**
- Protocol/agent yang query Miora reputation API bayar per-request via x402 (micropayment protocol dari Coinbase)
- Contoh: AI agent mau cek score wallet sebelum trade → kirim USDC micropayment via x402 → dapat score response
- Ini model yang proven — Cred Protocol sudah pakai x402 untuk monetize credit score API mereka

**Testnet:**
- EAS sudah deployed di Base Sepolia: bisa create schema, buat attestation, query attestation
- Gas di Base Sepolia gratis (pakai faucet)
- Bisa demo end-to-end: analyze wallet → publish score → read score on-chain

---

### Layer 2: Smart Follow + AI Alerts (Existing)

**Apa ini:**
User bisa follow wallet yang score-nya bagus dan mendapat notifikasi real-time ketika wallet tersebut trade, lengkap dengan AI risk assessment.

**Ini sudah di-build di V1:**
- 3-tier recommendation: Full Follow (80-100), Conditional Follow (40-79), Avoid (<40)
- Conditional follow dengan dynamic thresholds (computed dari data wallet sendiri)
- WebSocket real-time notifications
- AI risk assessment per notification via Gemini
- Email notifications via Resend
- Notification history di database

**Apa yang berubah di V2:**
- Recommendation sekarang juga di-publish on-chain sebagai bagian dari attestation
- User bisa lihat on-chain proof bahwa wallet ini "Full Follow" menurut Miora
- Tidak ada perubahan besar di layer ini — sudah solid

**Testnet:**
- Layer ini off-chain (backend + WebSocket), tidak perlu testnet
- Sudah functional, tinggal connect frontend

---

### Layer 3: AI Trading Agent (New)

**Apa ini:**
AI agent yang belajar dari wallet terbaik di Base dan trade otomatis untuk user — dengan budget, risk rules, dan conditions yang user tentukan sendiri.

**Bagaimana cara kerjanya:**
1. User connect wallet dan set parameters:
   - Budget: "Saya mau invest $500"
   - Max per trade: "$50"
   - Conditions: "Hanya token dengan liquidity > $100k, pair age > 6 jam, market cap > $1M"
   - Risk tolerance: "Low / Medium / High"
2. Miora AI agent berjalan di background:
   - Monitor top-scored wallets di Base (pakai scoring engine yang sudah ada)
   - Detect ketika top wallet melakukan trade
   - Evaluate trade pakai AI risk assessment (sudah ada)
   - Check apakah trade memenuhi user's conditions (sudah ada)
   - Check apakah masih dalam budget
3. Kalau semua pass → agent execute swap via Coinbase AgentKit + Agentic Wallet
4. User dapat notification: "Agent bought $30 of TOKEN_X because Wallet_ABC (score: 87) bought it. Liquidity $250k, pair age 2 days. Risk: Low."
5. User bisa pause, adjust rules, atau stop agent kapan saja

**Mengapa AgentKit + Agentic Wallets:**
- Dibuat oleh Coinbase, untuk Base — ini official infrastructure
- Agent punya wallet sendiri (Agentic Wallet) — user deposit budget ke wallet agent, agent manage sendiri
- Private key tidak pernah exposed ke agent — security by design
- Support gasless trading di Base
- Programmable guardrails — bisa set max spend per transaction, daily limits, dll
- Sudah ada SDK (Python dan TypeScript)

**Bedanya dengan copy-trade bot biasa:**
| Copy-trade bot biasa | Miora AI Agent |
|---|---|
| Blindly copy semua trade | Filter pakai scoring engine (6 faktor) |
| Tidak ada risk assessment | AI evaluate setiap trade sebelum execute |
| Tidak ada conditions | User set custom conditions (liquidity, pair age, mcap) |
| Tidak ada budget management | Budget limits, max per trade, stop-loss |
| Tidak tahu wallet mana yang bagus | Hanya follow wallet dengan score tinggi |

**Testnet:**
- AgentKit support Base Sepolia testnet
- Agentic Wallets bisa dibuat di testnet
- Swap bisa disimulasikan (atau pakai testnet DEX)
- Bisa demo: "Agent detected trade → evaluated risk → executed swap on Base Sepolia"

---

## Mengapa Konsep Ini Cocok untuk Base Batches

### 1. Alignment dengan Base 2026 Strategy

Base 2026 strategy punya 3 pilar. Miora V2 align dengan semuanya:

| Base 2026 Pilar | Miora V2 Alignment |
|---|---|
| "Building global markets — discover what's trending from top traders, grow your wealth" | Miora literally does this: discover top traders (scoring), follow them (watchlist), grow wealth (AI agent) |
| "Scaling payments and stablecoins" | x402 micropayments untuk reputation API, stablecoin-based agent trading |
| "Home for builders — agent-native smart accounts, x402, MCP" | Miora pakai AgentKit, Agentic Wallets, x402 — tiga tech stack utama Base |

### 2. Mengisi Gap di BB003 Cohort

12 project yang terpilih di BB003 Startup Track:
- AI agent infra: Blockrun, Agently, 4Mica, Floe Labs
- Lending: Credifi, Tomorrow
- Privacy: OPAL
- Prediction: Onsight, JPEG App, (Stealth)
- Neobank: Liminal
- FX: Nivo

**Yang tidak ada: trading intelligence / wallet reputation.** Miora mengisi gap ini.

### 3. Composable dengan Project BB003 Lain

Miora bukan isolated app — protocol lain bisa build di atasnya:
- **Credifi** (uncollateralized lending) bisa query Miora score untuk assess borrower trading quality
- **Agently** (agent marketplace) bisa integrate Miora sebagai "trading intelligence skill" untuk agents
- **Blockrun** (agent infra) bisa route agents ke Miora API untuk wallet analysis
- **Floe Labs** (agent credit) bisa pakai Miora score sebagai signal untuk agent creditworthiness

Ini bikin Miora jadi **connective tissue** di Base ecosystem, bukan standalone tool.

### 4. Deep Base Tech Integration

Miora V2 pakai 3 tech stack utama Base dalam satu produk:

| Tech | Dari | Fungsi di Miora |
|---|---|---|
| **EAS** (Ethereum Attestation Service) | Ethereum Foundation, live di Base | Publish trading reputation score on-chain |
| **AgentKit + Agentic Wallets** | Coinbase | AI agent autonomous trading |
| **x402** | Coinbase + Cloudflare | Monetize reputation API via micropayments |

Tidak ada project lain yang combine ketiganya. Juri akan lihat: "Dia pakai semua infrastructure kita secara kohesif."

### 5. Bukan "AI Wrapper" — Ada Proprietary Tech

Banyak hackathon project cuma pasang LLM di atas data dan bilang "AI-powered." Miora punya proprietary scoring engine:
- **FIFO buy-sell matching** untuk PnL calculation — ini non-trivial
- **6-factor scoring model** dengan configurable thresholds
- **Dynamic conditional thresholds** computed dari data wallet sendiri
- **3-tier recommendation system** dengan nuance (bukan binary)

AI (Gemini) dipakai untuk narasi dan risk assessment — tapi "otak" sebenarnya adalah scoring engine yang deterministic dan auditable.

---

## Competitive Landscape

### Wallet Reputation Space

| Project | Fokus | Chain | Perbedaan dengan Miora |
|---|---|---|---|
| Cred Protocol | Credit/lending risk scoring | Multi-chain | Fokus "apakah wallet ini bayar hutang?" — bukan trading quality |
| ChainAware | Fraud detection + wallet ranking | Multi-chain | Fokus security/fraud — bukan trading performance |
| Nomis | General reputation score | Multi-chain | Broad score, tidak spesifik trading |
| zScore | AI-based reputation | Multi-chain | Academic/research, belum production-ready |
| **Miora** | **Trading reputation** | **Base-first** | **Fokus "apakah wallet ini trader yang bagus?" — niche yang belum di-own** |

**Key insight:** Semua competitor fokus di credit risk atau fraud detection. Tidak ada yang fokus di **trading quality reputation**. Miora own niche ini.

### AI Trading Agent Space

| Project | Approach | Perbedaan dengan Miora |
|---|---|---|
| Generic copy-trade bots | Blindly copy all trades | Tidak ada scoring, tidak ada risk assessment, tidak ada conditions |
| GMGN | Wallet tracking + alerts | Alerts only, tidak ada autonomous execution |
| Cielo | Multi-chain wallet tracking | Analytics only, tidak ada agent |
| **Miora** | **Score-based intelligent agent** | **Agent hanya trade kalau wallet score tinggi + conditions met + AI risk assessment pass** |

---

## Technical Architecture V2

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                         │
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
│  │ [EXISTING]   │  │ [EXISTING]   │  │ [NEW]             │  │
│  └──────┬───────┘  └──────────────┘  └─────────┬─────────┘  │
│         │                                       │            │
│  ┌──────▼───────────────────────────────────────▼─────────┐  │
│  │              On-chain Layer (Base)                       │  │
│  │                                                         │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐  │  │
│  │  │ EAS         │  │ Agentic      │  │ x402          │  │  │
│  │  │ Attestation │  │ Wallet       │  │ Payments      │  │  │
│  │  │ [NEW]       │  │ [NEW]        │  │ [NEW]         │  │  │
│  │  └─────────────┘  └──────────────┘  └───────────────┘  │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              External APIs [EXISTING]                    │  │
│  │  Alchemy · DexScreener · Moralis · Gemini AI            │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Database (PostgreSQL) [EXISTING]            │  │
│  │  Users · Wallets · Transactions · Metrics ·             │  │
│  │  Watchlist · Notifications · Agent Configs              │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## Apa yang Berubah dari Codebase V1

| Komponen | Status | Detail |
|---|---|---|
| Scoring engine (6 faktor, FIFO PnL) | ✅ Tidak berubah | Core logic tetap sama |
| Conditional follow system | ✅ Tidak berubah | Dynamic thresholds tetap |
| AI risk assessment (Gemini) | ✅ Tidak berubah | Per-notification assessment tetap |
| Wallet monitoring service | ✅ Tidak berubah | Background polling tetap |
| WebSocket notifications | ✅ Tidak berubah | Real-time alerts tetap |
| Email notifications (Resend) | ✅ Tidak berubah | Async dispatch tetap |
| API clients (Alchemy, DexScreener, Moralis, Birdeye) | ✅ Tidak berubah | Data fetching tetap |
| Clean architecture (handlers, services, repos) | ✅ Tidak berubah | Pattern tetap |
| **EAS attestation integration** | 🆕 Baru | Publish score on-chain di Base |
| **x402 payment endpoint** | 🆕 Baru | Monetize reputation API |
| **AgentKit integration** | 🆕 Baru | Autonomous trading execution |
| **Agent config management** | 🆕 Baru | Budget, rules, pause/resume |
| **Agent dashboard UI** | 🆕 Baru | Monitor agent actions |
| **Multi-chain support** | ⬇️ De-prioritize | Fokus Base untuk hackathon, multi-chain jadi roadmap |

**Estimasi: 70% existing code tetap, 30% new code.**

---

## Testnet Deployment Plan

Semua komponen baru bisa di-test di Base Sepolia testnet:

### EAS Attestation
- EAS contracts sudah deployed di Base Sepolia
- Address: tersedia di docs EAS (attest.org)
- Flow: create schema → make attestation → query attestation
- SDK: `@ethereum-attestation-service/eas-sdk` (TypeScript) atau direct contract call dari Go
- Gas: gratis (Base Sepolia faucet)

### AgentKit + Agentic Wallet
- AgentKit SDK tersedia di npm (`@coinbase/agentkit`) dan PyPI (`coinbase-agentkit`)
- Support Base Sepolia untuk testing
- Agentic Wallet bisa dibuat di testnet
- Swap bisa disimulasikan atau pakai testnet liquidity

### x402 Payments
- x402 protocol support Base Sepolia
- Bisa setup payment endpoint yang terima USDC testnet
- Demo: agent kirim micropayment → dapat reputation score response

### Demo Flow di Testnet
```
1. User input wallet address di frontend
2. Backend analyze wallet, generate score
3. Score di-publish ke Base Sepolia via EAS → txn hash visible di BaseScan
4. User follow wallet dengan conditions
5. Wallet trade detected → AI assess risk → notification sent
6. AI agent evaluate → conditions met → execute swap di Base Sepolia
7. User lihat di dashboard: "Agent bought TOKEN_X, txn: 0x..."
```

Semua verifiable on-chain di Base Sepolia explorer.

---

## Pitch Script (30 detik)

> "Every day, thousands of wallets trade on Base. Some are great traders. Most aren't. But there's no way to tell the difference.
>
> Miora is a trading reputation protocol for Base. We analyze any wallet's trading history, compute a multi-factor score, and publish it on-chain as an EAS attestation — readable by any protocol, any agent, anyone.
>
> For users: follow the best wallets, get AI-filtered alerts, and let our agent trade for you with your rules.
>
> For protocols: query our reputation API to assess any wallet's trading quality.
>
> Built on Base's own infrastructure: EAS, AgentKit, x402. Base is home."

---

## Roadmap

### Phase 1: Hackathon (Now → April 27)
- [ ] Connect frontend ke backend (ganti dummy data)
- [ ] Deploy EAS schema + attestation di Base Sepolia
- [ ] Basic AgentKit proof of concept (detect trade → execute swap di testnet)
- [ ] x402 endpoint untuk reputation query
- [ ] Record founder video
- [ ] Submit application

### Phase 2: Post-Hackathon (Jika terpilih)
- [ ] Deploy ke Base mainnet
- [ ] Full AgentKit integration dengan budget management
- [ ] Agent dashboard UI
- [ ] Partnership dengan lending protocols untuk integrate reputation score
- [ ] Beta testing dengan real users

### Phase 3: Scale
- [ ] Multi-chain expansion (Ethereum, Arbitrum, Optimism, Polygon)
- [ ] Reputation score marketplace (protocols subscribe)
- [ ] Advanced agent strategies (multi-wallet, portfolio-based)
- [ ] Mobile app

---

## Summary

| Aspek | V1 (Sekarang) | V2 (Baru) |
|---|---|---|
| Kategori | Analytics tool | Infrastructure protocol + AI agent |
| On-chain presence | Tidak ada | EAS attestation, Agentic Wallet, x402 |
| Base-specific | Tidak ada | EAS + AgentKit + x402 (semua Base-native) |
| Monetization | Swap fee (belum jadi) | x402 micropayments + swap fee |
| Composability | Standalone app | Protocol lain bisa build di atas Miora |
| Competitive moat | Scoring engine | Scoring engine + on-chain reputation + agent intelligence |
| Narrative | "AI wallet analyzer" | "Trading reputation protocol for Base" |
| Wow factor | Iterasi dari yang sudah ada | Belum ada yang own "trading reputation" di Base |
