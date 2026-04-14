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
| 6 | What part is onchain? | ✅ Done | Akan lebih kuat kalau fee router sudah deploy |
| 7 | Ideal customer profile | ✅ Done | |
| 8 | Category | ✅ Done | |
| 9 | Location | ❌ Kamu isi | Personal info |
| 10 | Token? | ✅ Done | |
| 11 | What part uses Base? | ✅ Done | |
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
| 26 | Dune / contract addresses? | ❌ Kamu isi | Deploy fee router dulu di Base Sepolia |
| 27 | Why join Base Batches? | ✅ Done | |
| 28 | Anything else? | ❌ Kamu isi | Optional tapi recommended |
| 29 | Who referred you? | ❌ Kamu isi | Kalau ada |
| 30 | GitHub repo | ✅ Done | |
| 31 | Private repo collaborator? | ✅ Done | |

**Summary: 15 sudah dijawab ✅ — 16 perlu kamu isi ❌ (mayoritas personal info)**

### ⚠️ Jawaban yang Sudah Saya Tulis Tapi Perlu Kamu Review

| # | Pertanyaan | Kenapa Ragu |
|---|---|---|
| 6 | What part is onchain? | Saat ini jawaban terasa tipis — mayoritas "baca data on-chain", bukan "logic on-chain." Fee router belum deploy. Kalau juri strict soal on-chain component, ini bisa jadi kelemahan. |
| 11 | What part uses Base? | Saya tulis "Base is the primary chain" tapi kenyataannya di codebase tidak ada yang Base-exclusive. Semua fitur jalan di 5 EVM chain. Juri bisa challenge ini. Kamu perlu decide apakah mau reframe atau jujur saja. |
| 17 | How far along? | Saya pilih "Prototype" — tapi bisa juga argue "MVP" karena backend fully functional. Tergantung bagaimana kamu define. Prototype lebih honest karena frontend belum connect. |
| 20 | What part is magic? | Jawaban saya fokus ke scoring engine. Tapi ini subjektif — mungkin kamu merasa bagian lain yang lebih "magic" (misalnya AI risk assessment per notification, atau conditional follow system). Kamu yang paling tahu apa yang paling kamu banggakan. |
| 21 | Unique insight? | Saya tulis "90% users don't need more charts, they need decisions." Ini catchy tapi belum di-validate dengan data. Kalau juri tanya "based on what?", kamu perlu punya jawaban. Idealnya ada personal story atau user research yang back this up. |
| 23-25 | Users/Revenue | Saya jawab "not yet / N/A" — ini honest tapi lemah. Kalau sebelum submit kamu bisa dapat bahkan 5-10 orang coba produknya dan kasih feedback, jawaban ini berubah drastis. |
| 27 | Why join Base Batches? | Jawaban saya generic — mentorship, ecosystem, Demo Day. Setiap applicant akan jawab hal serupa. Akan lebih kuat kalau kamu tambahkan alasan spesifik yang personal (misalnya: "I've been trading on Base for X months and noticed Y problem firsthand"). |

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

AI-powered wallet intelligence and trading on Base

---

## What is your product's unique value proposition?

Miora AI turns overwhelming on-chain data into simple, actionable decisions. Instead of reading charts and numbers, users get a clear answer: "Follow this wallet", "Follow with conditions", or "Avoid." 

Existing tools like Nansen and Arkham are built for power users. Miora is built for everyone — analyze any wallet, get an AI-scored recommendation, set smart alert conditions, and trade directly from the platform. No more switching between analytics dashboards, Telegram alpha groups, and DEX frontends.

---

## What part of your product is onchain?

- **Wallet analysis** — Reads on-chain transaction history from Base via Alchemy to compute scoring metrics (PnL, win rate, entry timing, trade discipline)
- **Token swap execution** — Executes swaps on Base through 1inch aggregator (Uniswap, SushiSwap, Curve, Balancer routing)
- **Market data enrichment** — Pulls on-chain pair data (liquidity, market cap, pair age) from DexScreener and Moralis for Base tokens
- **Fee Router smart contract** — [PLANNED] Solidity contract on Base for collecting swap fees as the monetization layer

---

## What is your ideal customer profile?

Retail crypto traders on Base who want to discover profitable wallets to follow and trade smarter — without needing advanced analytics skills. Specifically:

- Traders who moved from CEX to DEX and feel overwhelmed by raw on-chain data
- Users who see "whale alerts" on Twitter but don't know how to evaluate if a wallet is actually good
- Beginners who want guidance on who to follow and when to pay attention

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

Miora supports wallet analysis, token swaps, and smart alerts on Base chain. Specifically:

- **Wallet analysis on Base** — Fetch transaction history via Alchemy, calculate PnL with FIFO buy-sell matching, generate multi-factor scoring
- **Swap quotes and execution on Base** — Route through 1inch to find best prices across Base DEXs (Uniswap, Aerodrome, SushiSwap, etc.)
- **Historical price data on Base** — Moralis for block-level token prices on Base
- **Real-time pair data on Base** — DexScreener for liquidity, market cap, pair age of Base tokens
- **Smart alerts for Base wallets** — Monitor followed wallets on Base, notify when they trade with AI risk assessment

Base is the primary chain. Multi-chain support (Ethereum, Arbitrum, Optimism, Polygon, Solana) exists as an expansion layer, but Base is the focus for go-to-market.

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

[FILL — e.g., "I (Yobel) handle all technical development — backend (Go), frontend (Next.js/TypeScript), smart contracts (Solidity/Anchor), and infrastructure."]

---

## How long have the founders known each other and how did you meet?

[FILL — If solo: "Solo founder."]

---

## How far along are you?

Prototype

Backend is complete (Go + Fiber, clean architecture, all API endpoints functional). Frontend is complete (Next.js 16, all pages built). Currently connecting frontend to backend API. Smart contract fee router is in development.

---

## How long have you been working on this?

[FILL]

---

## How much of that time is full-time vs part-time?

[FILL]

---

## What part of your product is magic or impressive?

The scoring engine. Miora doesn't just show you wallet data — it computes a multi-factor score using:

1. **FIFO buy-sell matching** — Accurately calculates realized and unrealized PnL by matching buys to sells in order, not just comparing balances
2. **Six scoring factors** — Win rate, profit consistency (standard deviation), entry timing (how early into new tokens), token quality (log-scale market cap), trade discipline (focus ratio), and risk exposure (low-liquidity token percentage)
3. **Three-tier recommendation** — Score 80-100: Full Follow (safe to copy all trades). Score 40-79: Conditional Follow with AI-generated filter conditions (e.g., "only notify if liquidity > $100k and pair age > 6 hours"). Score < 40: Avoid.
4. **Dynamic condition thresholds** — Conditions are computed from the wallet's own trading data (median liquidity, market cap, volume, average pair age), not hardcoded numbers
5. **AI risk assessment per alert** — When a followed wallet trades, Gemini evaluates the token's market data and gives a plain-language risk opinion before you decide to act

No other tool gives you a simple "should I follow this wallet?" answer with customizable, intelligent alert filters.

---

## What is your unique insight or advantage?

90% of crypto users don't need more charts — they need someone to tell them which wallets are worth following and when to pay attention.

Existing wallet analytics tools (Nansen, Arkham, DeBank) show data. Miora shows decisions. The gap in the market isn't "better analytics" — it's "actionable intelligence for non-experts." We combine wallet scoring, conditional alerts, and DEX trading into one flow so users go from discovery to action without leaving the platform.

The technical advantage: our scoring engine uses FIFO PnL matching and multi-factor analysis that goes beyond simple profit tracking. Most "wallet score" tools just look at total profit. We evaluate how consistently they profit, how early they enter, what quality tokens they trade, and how disciplined their strategy is.

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

N/A — pre-launch. Planned revenue model: swap fee collection via on-chain fee router contract on Base.

---

## Dune dashboards / deployed contract addresses

[FILL — Deploy fee router on Base Sepolia and paste address here before submitting]

---

## Why do you want to join Base Batches?

Three reasons:

1. **Mentorship and feedback** — As a solo undergraduate founder, direct access to experienced builders and the Base team would accelerate product decisions that currently take me weeks of research
2. **Base ecosystem alignment** — Miora's mission of making on-chain trading accessible aligns with Base's goal of bringing the next wave of users onchain. Being embedded in the Base ecosystem means better integration support and distribution
3. **Demo Day exposure** — Presenting to investors and the Base community in San Francisco would be a transformative opportunity to validate the product direction and explore funding for full-time development after graduation

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
- [ ] Connect frontend to backend API (replace dummy data)
- [ ] Deploy fee router contract on Base Sepolia
- [ ] Fill all [FILL] fields
- [ ] Make GitHub repo public
- [ ] Review all answers for clarity and conciseness
- [ ] Submit before April 27, 2026
