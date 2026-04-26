# Base Batches 003: Student Track — Application Draft

---

## 📊 Status Overview

| # | Pertanyaan | Status | Catatan |
|---|---|---|---|
| 1 | Company Name | ✅ Done | |
| 2 | Website / Product URL | ❌ Kamu isi | Butuh deploy URL atau link |
| 3 | Demo URL | ❌ Kamu isi | Butuh video/prototype link |
| 4 | Describe what your company does (~50 chars) | ✅ Done | |
| 5 | Unique value proposition | ✅ Done | |
| 6 | What part is onchain? | ✅ Done | EAS attestation + AgentKit bots + auto profit transfer |
| 7 | Ideal customer profile | ✅ Done | |
| 8 | Category | ✅ Done | |
| 9 | Location | ❌ Kamu isi | Personal info |
| 10 | Token? | ✅ Done | |
| 11 | What part uses Base? | ✅ Done | Base-exclusive — EAS + AgentKit + wallet auth |
| 12 | Founder names + contact | ❌ Kamu isi | Personal info |
| 13 | Founder background + LinkedIn | ❌ Kamu isi | Personal info |
| 14 | 1-min founder video | ❌ Kamu isi | Harus direkam sendiri — ini sangat penting |
| 15 | Who writes code? | ❌ Kamu isi | Saya kasih contoh, tapi kamu yang confirm |
| 16 | How long founders known each other? | ❌ Kamu isi | Solo atau tim |
| 17 | How far along? | ✅ Done | |
| 18 | How long working on this? | ❌ Kamu isi | Hanya kamu yang tahu |
| 19 | Full-time vs part-time? | ❌ Kamu isi | Hanya kamu yang tahu |
| 20 | What part is magic/impressive? | ✅ Done | |
| 21 | Unique insight/advantage? | ✅ Done | |
| 22 | Plan to raise VC / launch token? | ❌ Kamu isi | Keputusan personal |
| 23 | Users/customers? | ✅ Done | Jawaban: belum ada |
| 24 | Active users / paying? | ✅ Done | N/A |
| 25 | Revenue? | ✅ Done | N/A |
| 26 | Dune / contract addresses? | ❌ Kamu isi | Paste EAS schema UID + attestation explorer link |
| 27 | Why join Base Batches? | ✅ Done | |
| 28 | Anything else? | ❌ Kamu isi | Optional tapi recommended |
| 29 | Who referred you? | ❌ Kamu isi | Kalau ada |
| 30 | GitHub repo | ✅ Done | |
| 31 | Private repo collaborator? | ✅ Done | |

**Summary: 15 sudah dijawab ✅ — 16 perlu kamu isi ❌ (mayoritas personal info)**

### ⚠️ Jawaban yang Sudah Saya Tulis Tapi Perlu Kamu Review

| # | Pertanyaan | Kenapa Perlu Review |
|---|---|---|
| 17 | How far along? | Saya tulis "Prototype" — tapi bisa argue "MVP" karena backend + frontend + agent sidecar semua sudah built. EAS schema sudah registered di Base Sepolia. Prototype lebih honest karena frontend belum fully connected ke backend. |
| 20 | What part is magic? | Jawaban fokus ke scoring engine + bot system. Tapi ini subjektif — mungkin kamu merasa bagian lain yang lebih "magic" (misalnya AI risk assessment per notification, atau conditional follow system). |
| 21 | Unique insight? | Saya tulis "90% users don't need more charts, they need decisions." Catchy tapi belum di-validate dengan data. Kalau juri tanya "based on what?", kamu perlu punya jawaban. Idealnya ada personal story atau user research. |
| 23-25 | Users/Revenue | Saya jawab "not yet / N/A" — honest tapi lemah. Kalau sebelum submit kamu bisa dapat 5-10 orang coba produknya, jawaban ini berubah drastis. |
| 27 | Why join Base Batches? | Jawaban agak generic. Akan lebih kuat kalau kamu tambahkan alasan spesifik yang personal (misalnya: "I've been trading on Base for X months and noticed Y problem firsthand"). |

---

> ⚠️ Fields marked [FILL] need to be completed manually.

---

## Company Name

Miora AI

---

## Website / Product URL

[FILL — deploy URL or GitHub link]

---

## Demo URL

[FILL — video recording or live demo link]

---

## Describe what your company does (~50 chars)

Trading reputation protocol + AI bots on Base

---

## What is your product's unique value proposition?

Miora AI turns overwhelming on-chain data into simple, actionable decisions. Instead of reading charts and numbers, users get a clear answer: "Follow this wallet" (score 80+), "Follow with conditions" (score 40-79), or "Avoid" (score <40).

What makes Miora different from existing analytics tools:

1. **Decisions, not data** — Nansen and Arkham show raw data for power users. Miora shows a recommendation anyone can act on.
2. **On-chain reputation** — Every score is published as an EAS attestation on Base. Other protocols, agents, and dApps can query and use these scores — it's a composable primitive, not just a number in our database.
3. **Intelligent automation** — Two bot types powered by Coinbase AgentKit: Wallet Bot copies one trader's moves, Consensus Bot trades when multiple top wallets agree. Both evaluate conditions and AI risk assessment before every trade.
4. **Auto profit transfer** — Sell proceeds are automatically transferred to the user's connected wallet. No manual withdrawal needed.
5. **Wallet-based auth** — MetaMask connect only. No email, no password, no centralized login. Pure Web3.

Users go from discovery to action without leaving the platform: analyze → follow → automate.

---

## What part of your product is onchain?

- **EAS Attestation** — Trading reputation scores are published on-chain on Base Sepolia via Ethereum Attestation Service. Each analyzed wallet gets an attestation with score, recommendation, total transactions, and chain — verifiable by any protocol or agent. Schema: `uint8 score, string recommendation, uint32 totalTransactions, string chain`.
- **Wallet analysis** — Reads on-chain transaction history from Base via Alchemy to compute scoring metrics (FIFO PnL, win rate, entry timing, trade discipline, token quality, risk exposure).
- **AI Trading Bots** — Two bot types powered by Coinbase AgentKit + Agentic Wallets:
  - Wallet Bot: copies one wallet's buys AND sells on Base
  - Consensus Bot: trades when 3+ high-score wallets buy the same token within a time window
  - All bot transactions are on-chain and visible on BaseScan
- **Auto profit transfer** — When a bot sells, proceeds are automatically transferred on-chain from the Agentic Wallet to the user's connected wallet.
- **Market data enrichment** — Pulls on-chain pair data (liquidity, market cap, pair age) from DexScreener and historical token prices from Moralis for Base tokens.

---

## What is your ideal customer profile?

Retail crypto traders on Base who want to discover profitable wallets to follow and trade smarter — without needing advanced analytics skills. Specifically:

- Traders who moved from CEX to DEX and feel overwhelmed by raw on-chain data
- Users who see "whale alerts" on Twitter but can't evaluate if a wallet is actually good
- Beginners who want guidance on who to follow and when to pay attention
- Traders who want to automate copy-trading with intelligent conditions, not blind copying

---

## Which category best describes your company?

DeFi / AI

---

## Where are you located now, and where would the company be based after the program?

[FILL]

---

## Do you already have a token?

No.

---

## What part of your product uses Base?

Miora is built exclusively on Base. Every on-chain component runs on Base Sepolia:

- **EAS Attestation on Base** — Trading reputation scores published as on-chain attestations via EAS (contract: `0x4200000000000000000000000000000000000021`). Verifiable on base-sepolia.easscan.org.
- **Wallet analysis on Base** — Fetch transaction history via Alchemy Base RPC, calculate PnL with FIFO buy-sell matching, generate multi-factor scoring (6 factors).
- **AI Trading Bots on Base** — Two bot types via Coinbase AgentKit + Agentic Wallets on Base Sepolia:
  - Wallet Bot: copies one wallet's buys AND sells
  - Consensus Bot: trades when multiple high-score wallets buy the same token
  - Auto profit transfer: sell proceeds sent on-chain to user's connected wallet
- **Historical price data on Base** — Moralis for block-level token prices on Base.
- **Real-time pair data on Base** — DexScreener for liquidity, market cap, pair age of Base tokens.
- **Smart alerts for Base wallets** — Monitor followed wallets on Base, notify when they trade with AI risk assessment (Gemini).
- **Wallet-based auth** — MetaMask connect via wagmi/viem. No Firebase, no centralized auth. `X-Wallet-Address` header for API authentication.

Built on Base's own infrastructure: EAS + AgentKit. Base is home.

---

## Founder(s) Names and Contact Information

[FILL]

---

## Please describe each founder's background and add their LinkedIn profile(s)

[FILL]

---

## Please enter the URL of a ~1-minute unlisted video introducing the founder(s) and what you're building

[FILL — Record a 1-minute video covering: who you are, the problem, what Miora does, why Base]

---

## Who writes code or handles technical development?

[FILL — e.g., "I (Yobel) handle all technical development — backend (Go + Fiber), frontend (Next.js + TypeScript), agent sidecar (Python + AgentKit), and infrastructure (Docker)."]

---

## How long have the founders known each other and how did you meet?

[FILL — If solo: "Solo founder."]

---

## How far along are you?

Prototype

Full stack is built and functional:
- **Backend** (Go + Fiber): All API endpoints working — wallet analysis, scoring engine, watchlist CRUD, bot management, EAS attestation, WebSocket notifications. Clean architecture with DI container.
- **Frontend** (Next.js 16 + Tailwind v4 + shadcn/ui): All pages built — landing, analyze, watchlist, bot management, bot detail. Wallet-based auth via MetaMask (Reown AppKit).
- **Agent sidecar** (Python + FastAPI + Coinbase AgentKit): Agentic Wallet creation, swap execution, balance queries, transfer. Running on port 8090.
- **On-chain**: EAS schema registered on Base Sepolia. Attestation flow integrated into wallet analysis.
- **Infrastructure**: Docker Compose with 4 services (PostgreSQL, backend, frontend, agent sidecar).

Currently connecting frontend to backend API (replacing dummy data with real API calls).

---

## How long have you been working on this?

[FILL]

---

## How much of that time is full-time vs part-time?

[FILL]

---

## What part of your product is magic or impressive?

Two things stand out:

**1. The Scoring Engine**

Miora doesn't just show wallet data — it computes a multi-factor score that produces an actionable recommendation:

- **FIFO PnL matching** — Accurately calculates realized and unrealized PnL by matching buys to sells in order, not just comparing balances
- **Six scoring factors** — Win rate, profit consistency (standard deviation), entry timing (how early into new tokens), token quality (log-scale market cap), trade discipline (focus ratio), and risk exposure (low-liquidity token percentage)
- **Three-tier recommendation** — Score 80-100: Full Follow (safe to copy all trades). Score 40-79: Conditional Follow with dynamic filter conditions. Score < 40: Avoid.
- **Dynamic condition thresholds** — Conditions are computed from the wallet's own trading data (median liquidity, market cap, volume, average pair age), not hardcoded numbers

**2. The Bot System**

Two bot types that go beyond blind copy-trading:

- **Wallet Bot** — Copies one wallet's buys AND sells. Conditions auto-filled from analyze result. Every trade is evaluated through conditions + AI risk assessment before execution.
- **Consensus Bot** — Scans all Miora-analyzed wallets. Trades only when 3+ high-score wallets buy the same token within a configurable time window. Higher confidence through crowd intelligence.
- **Auto profit transfer** — Sell proceeds are automatically transferred on-chain to the user's connected wallet. No manual withdrawal.

No other tool gives you a simple "should I follow this wallet?" answer with intelligent conditions AND autonomous bots that act on it.

---

## What is your unique insight or advantage?

90% of crypto users don't need more charts — they need someone to tell them which wallets are worth following and when to pay attention.

Existing wallet analytics tools (Nansen, Arkham, DeBank) show data. Miora shows decisions. The gap in the market isn't "better analytics" — it's "actionable intelligence for non-experts." We combine wallet scoring, conditional alerts, and AI trading bots into one flow so users go from discovery to action without leaving the platform.

**Technical advantage:** Our scoring engine uses FIFO PnL matching and multi-factor analysis that goes beyond simple profit tracking. Most "wallet score" tools just look at total profit. We evaluate how consistently they profit, how early they enter, what quality tokens they trade, and how disciplined their strategy is.

**Composability advantage:** Every score is published on-chain via EAS. This means other protocols can query Miora scores — lending protocols can assess borrower quality, AI agents can check wallet reputation before copy-trading, dApps can gate features based on trading score. Miora becomes infrastructure, not just a tool.

---

## Do you plan on raising capital from VCs? Do you plan to launch a token?

[FILL]

---

## Do you have users or customers?

Not yet — currently in prototype stage. Plan to onboard beta testers on Base before public launch.

---

## Active users / paying customers

N/A — pre-launch.

---

## Revenue

N/A — pre-launch. Planned revenue model:
- **Consensus Bot** as premium feature (subscription or per-trade fee) — trades with higher confidence using multi-wallet agreement signals
- **Reputation API** at scale — public endpoint `GET /api/reputation/:address` is free now, monetizable via subscription or per-query fees when demand grows

---

## Dune dashboards / deployed contract addresses

[FILL — Paste EAS schema UID and attestation explorer link. Example: https://base-sepolia.easscan.org/schema/view/0x...]

---

## Why do you want to join Base Batches?

Three reasons:

1. **Mentorship and feedback** — As a solo undergraduate founder, direct access to experienced builders and the Base team would accelerate product decisions that currently take me weeks of research. Specifically, I want guidance on scaling the scoring engine and optimizing bot execution for mainnet.
2. **Base ecosystem alignment** — Miora is built exclusively on Base's own infrastructure (EAS + AgentKit). Being embedded in the Base ecosystem means better integration support, early access to new tools, and distribution to Base's 34M+ monthly active users.
3. **Demo Day exposure** — Presenting to investors and the Base community in San Francisco would be a transformative opportunity to validate the product direction and explore funding for full-time development after graduation.

---

## Anything else you'd like us to know?

[FILL — Optional: mention hackathon experience, relevant coursework, personal motivation, etc.]

---

## Who referred you to this program?

[FILL]

---

## GitHub repo link

https://github.com/yebology/miora-ai

---

## If your GitHub repo is private, have you added "devfolio-judge" as a collaborator?

Not Applicable (repo is public)

---

## 📋 Pre-Submission Checklist

- [ ] Record 1-minute founder video
- [x] Register EAS schema on Base Sepolia (`make register-schema`)
- [ ] Fill all [FILL] fields
- [ ] Paste EAS schema UID + attestation explorer link in Dune/contract addresses field
- [ ] Make GitHub repo public
- [ ] Review all answers for clarity and conciseness
- [ ] Submit before April 27, 2026
