# Miora AI — Pitch Deck

---

## Slide 1: Title

**Miora AI**
Score any trader on Base. Let AI ride the winners for you.

Trading Reputation Protocol + AI Trading Bots on Base

Base Batches 003 — Student Track

---

## Slide 2: Problem

Every day, thousands of wallets trade on Base. Some are great traders. Most aren't. But there's no way to tell the difference.

**The numbers:**
- 80-90% of retail crypto traders lose money over 6-24 months (Tapbit, 2026)
- 97% of persistent day traders (300+ days) end up losing money (CoinCub, 2026)
- Base processes 10-11M transactions daily with 34M+ monthly active users (BaseScan, 2025)
- $17T stablecoin volume on Base in 2025 (Base 2026 Strategy)

**The problem:**
- Retail traders see "whale alerts" on Twitter but can't evaluate if a wallet is actually good
- Existing tools (Nansen, Arkham) show raw data — charts, numbers, transaction lists
- Users don't need more data. They need decisions: "Should I follow this wallet?"
- Even if they find a good wallet, they can't monitor it 24/7

The gap: actionable trading intelligence for non-experts.

Sources: [Tapbit](https://blog.tapbit.com/why-80-90-of-crypto-day-traders-lose-money-and-how-to-be-in-the-10-20-that-dont-2026-guide/), [CoinCub](https://coincub.com/blog/crypto-trading-traps/), [Base 2026 Strategy](https://blog.base.org/2026-mission-vision-and-strategy), [BaseScan](https://basescan.org/)

---

## Slide 3: Solution

**You paste a wallet. We tell you if it's worth following.**

Miora turns any wallet address into a simple answer:
- "Follow this wallet" (score 80+)
- "Follow with conditions" (score 40-79)
- "Avoid" (score <40)

Then you choose what happens next:
- **Watch** — get notified when they trade, with AI risk assessment
- **Automate** — let a bot copy their trades with your budget and rules

Every score is published on-chain via EAS on Base — verifiable by anyone, queryable by any protocol.

---

## Slide 4: Under the Hood

What makes Miora's scoring different from "just another analytics tool":

- **FIFO PnL matching** — actual profit calculation, not just balance snapshots
- **6 scoring factors** — win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure
- **Dynamic conditions** — thresholds computed from the wallet's own data, not hardcoded
- **AI risk assessment** — Gemini generates a plain-language risk opinion for every trade and notification, recorded in trade history for user review
- **On-chain proof** — scores published as EAS attestations, not just a number in our database

---

## Slide 5: AI Trading Bots

Two bot types, powered by **Coinbase AgentKit + Agentic Wallets**:

**What is an Agentic Wallet?**
A wallet created and managed by Coinbase's AgentKit — the bot has its own wallet on Base to trade with. User deposits funds, bot manages them autonomously. Key management handled by Coinbase Developer Platform.

**Wallet Bot**
- Pick a wallet from your watchlist
- Bot copies its buys AND sells automatically
- Conditions auto-filled from analyze result
- Set budget + max per trade → bot handles the rest

**Consensus Bot (Premium)**
- Scans ALL wallets analyzed by Miora
- Trades only when 3+ high-score wallets buy the same token within a time window
- Higher confidence — crowd intelligence, not just one wallet
- Revenue stream: premium feature

---

## Slide 6: Why Base

Built exclusively on Base's own infrastructure:

- **EAS** (Ethereum Attestation Service) — on-chain reputation scores, same standard used by Coinbase Verifications
- **Coinbase AgentKit** — autonomous trading via Agentic Wallets
- Base is the #1 L2 by revenue, 46% of all L2 DeFi TVL (Sherlock, 2026)
- $17T stablecoin volume on Base in 2025 (Base 2026 Strategy)
- 10-11M daily transactions, 34M+ monthly active users (BaseScan/CoinLedger, 2025)

Base 2026 vision: "Discover what's trending from top traders, grow your wealth."
Miora literally does this.

No other project in BB003 cohort covers trading intelligence or wallet reputation.

Sources: [Sherlock](https://sherlock.xyz/post/best-blockchain-to-build-on-in-2026), [Base 2026 Strategy](https://blog.base.org/2026-mission-vision-and-strategy), [BB003 Blog](https://blog.base.org/introducing-base-batches-003-2), [CoinLedger](https://coinledger.io/nl/research/base-tvl-and-network-growth)

---

## Slide 7: Competitive Landscape

Nobody owns "trading quality reputation" on Base.

| | Nansen/Arkham | Copy-trade bots | Cred Protocol | Miora AI |
|---|---|---|---|---|
| Shows data | ✅ | ❌ | ❌ | ✅ |
| Shows decisions | ❌ | ❌ | ❌ | ✅ |
| On-chain reputation | ❌ | ❌ | ✅ (credit) | ✅ (trading) |
| AI risk assessment | ❌ | ❌ | ❌ | ✅ |
| Smart conditions | ❌ | ❌ | ❌ | ✅ |
| Autonomous trading | ❌ | ✅ (blind) | ❌ | ✅ (intelligent) |
| Base-native | ❌ | ❌ | Multi-chain | ✅ |

Competitors focus on credit risk (Cred Protocol), fraud detection (ChainAware), or general reputation (Nomis). Nobody focuses on **trading quality reputation**.

Sources: [Cred Protocol](https://credprotocol.com), [ChainAware](https://chainaware.ai), [Web3 Reputation Comparison 2026](https://chainaware.ai/blog/web3-reputation-score-comparison-2026/)

---

## Slide 8: Business Model

**Free tier**
- Analyze wallets (unlimited)
- Follow wallets + notifications
- Wallet bot (copy one wallet)

**Premium (Consensus Bot)**
- Scan all Miora wallets for consensus signals
- Trade when multiple high-score wallets agree
- Subscription or per-trade fee

**B2B (Reputation API) — Future**
- Lending protocols query Miora scores to assess borrowers
- AI agents check wallet reputation before copy-trading
- Public API endpoint — monetization via subscription or per-query fees

---

## Slide 9: Demo

**Live on Base Sepolia — every step verifiable on-chain.**

[Insert product screenshots or live demo video]

---

## Slide 10: Roadmap

**Now (Student Track)**
- Deploy EAS attestation on Base Sepolia
- Bot PoC: wallet bot + consensus bot
- Connect frontend to backend

**Post-Program**
- Base mainnet deployment
- Consensus bot as paid premium feature
- Full DEX integration (Aerodrome/Uniswap)
- Reputation leaderboard

**Miora AI — Score. Follow. Bot. On Base.**
