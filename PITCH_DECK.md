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
- **Watch** — get notified when they trade, with AI risk assessment per alert
- **Automate** — let a bot copy their trades with your budget and rules. Sell proceeds auto-transfer to your connected wallet.

Every score is published on-chain via EAS on Base — verifiable by anyone, queryable by any protocol.

---

## Slide 4: Under the Hood

What makes Miora's scoring different from "just another analytics tool":

- **FIFO PnL matching** — actual profit calculation, not just balance snapshots
- **6 scoring factors** — win rate, profit consistency, entry timing, token quality, trade discipline, risk exposure
- **Dynamic conditions** — thresholds computed from the wallet's own data (median liquidity, mcap, volume, pair age), not hardcoded
- **AI risk assessment** — Gemini generates a plain-language risk opinion for every trade notification, recorded in trade history for user review
- **On-chain proof** — scores published as EAS attestations on Base, not just a number in our database

---

## Slide 5: AI Trading Bots

Two bot types, powered by **Coinbase AgentKit + Agentic Wallets**:

**What is an Agentic Wallet?**
A wallet created and managed by Coinbase's AgentKit on Base. The bot has its own wallet to trade with. Key management handled by Coinbase Developer Platform — we never touch the private key.

**Wallet Bot**
- Pick a wallet from your watchlist
- Bot copies its buys AND sells automatically
- Conditions auto-filled from analyze result (e.g., min liquidity, min mcap)
- Set budget + max per trade → bot handles the rest
- Sell proceeds auto-transfer to your connected wallet

**Consensus Bot (Premium)**
- Scans ALL wallets analyzed by Miora
- Trades only when 3+ high-score wallets buy the same token within a time window
- Higher confidence — crowd intelligence, not just one wallet
- Configurable: min score, consensus threshold, time window

---

## Slide 6: Why Base

Built exclusively on Base's own infrastructure:

- **EAS** (Ethereum Attestation Service) — on-chain reputation scores, same standard used by Coinbase Verifications
- **Coinbase AgentKit** — autonomous trading via Agentic Wallets on Base
- **Wallet-based auth** — MetaMask connect via wagmi/viem, no centralized auth
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
| Auto profit transfer | ❌ | ❌ | ❌ | ✅ |
| Base-native | ❌ | ❌ | Multi-chain | ✅ |

Competitors focus on credit risk (Cred Protocol), fraud detection (ChainAware), or general reputation (Nomis). Nobody focuses on **trading quality reputation**.

Sources: [Cred Protocol](https://credprotocol.com), [ChainAware](https://chainaware.ai), [Web3 Reputation Comparison 2026](https://chainaware.ai/blog/web3-reputation-score-comparison-2026/)

---

## Slide 8: Business Model

**Free tier**
- Analyze wallets (unlimited)
- Follow wallets + real-time notifications with AI risk assessment
- Wallet bot (copy one wallet's trades, auto profit transfer)

**Premium (Consensus Bot)**
- Scan all Miora wallets for consensus signals
- Trade when multiple high-score wallets agree
- Subscription or per-trade fee

**B2B (Reputation API) — Future**
- Public reputation endpoint: `GET /api/reputation/:address`
- Lending protocols query Miora scores to assess borrowers
- AI agents check wallet reputation before copy-trading
- Monetization via subscription or per-query fees at scale

---

## Slide 9: Demo

**Live on Base Sepolia — every step verifiable on-chain.**

[Insert product screenshots or live demo video]

---

## Slide 10: Roadmap

**Now (Student Track)**
- ✅ EAS attestation deployed on Base Sepolia
- ✅ Two bot types: wallet bot + consensus bot (AgentKit)
- ✅ Full backend + frontend built
- 🔄 Connecting frontend to backend API

**Post-Program**
- Base mainnet deployment
- Consensus bot as paid premium feature
- Full DEX integration for bot swaps (Aerodrome/Uniswap)
- Reputation leaderboard
- Multi-bot portfolio tracking

**Miora AI — Score. Follow. Bot. On Base.**
