# Miora AI — Research & Sources

Dokumen ini berisi semua hasil riset yang mendasari pivot Miora V1 → V2 dan keputusan untuk submit ke Base Batches Student Track.

---

## 1. Riset Step-by-Step (Kronologis)

### Step 1: Scan Codebase Miora V1
- Baca README.md, PROGRESS.md, semua backend services, entities, handlers, clients, config, constants, router
- Baca semua frontend pages, components, constants, types, hooks
- **Temuan**: Backend 100% done (clean architecture, scoring engine, FIFO PnL, conditional follow, AI insights, WebSocket, email alerts). Frontend done tapi pakai dummy data. Smart contracts belum jadi. Multi-chain (5 EVM + Solana).

### Step 2: Evaluasi Konsep V1 untuk Solana Frontier
- Penilaian objektif: 6.5/10 secara konsep
- Kelemahan utama: bukan Solana-native (multi-chain generic), analytics tool sudah saturated (Nansen, Arkham, DeBank, Cielo), AI wrapper bukan differentiator di 2026
- Solana integration tipis: cuma API calls (Jupiter, Birdeye, Alchemy Solana)

### Step 3: Evaluasi untuk Base Batches Student Track
- Baca halaman Base Batches Student Track (devfolio)
- Baca blog "Introducing Base Batches 003" — 12 tim terpilih dari 1,175 applicants
- Baca Base 2026 Mission, Vision, Strategy
- **Temuan**: Miora secara DNA adalah EVM product (5/6 chain = EVM termasuk Base). Base Batches = application-based bukan hackathon ("It doesn't matter if you had previously written code"). Student track deadline April 27.

### Step 4: Perbandingan Solana Frontier vs Base Batches
- Scoring: Solana fit 5.5/10, Base fit 7.5/10
- Alasan Base lebih cocok: chain alignment (EVM-native), format application-based (belum selesai bukan masalah), narrative fit ("beginner-friendly" = Base mission), kompetisi lebih favorable (student track)

### Step 5: Analisis BB003 Cohort
- Riset 12 tim yang terpilih: Blockrun.ai, Stealth, 4Mica, OPAL, Onsight, Credifi, Tomorrow, Agently, Nivo, JPEG App, Floe Labs, Liminal
- **Temuan**: 5/12 = AI agent projects. Gap: tidak ada trading intelligence / wallet reputation. Miora bisa isi gap ini.

### Step 6: Riset Base Ecosystem Trends
- Base 2026 strategy: 3 pilar (global markets, stablecoin payments, home for builders/agents)
- EAS (Ethereum Attestation Service) sudah live di Base
- Coinbase AgentKit + Agentic Wallets — agent infrastructure
- x402 protocol — micropayments
- $17T stablecoin volume di Base 2025

### Step 7: Konsep V2 — Trading Reputation Protocol
- Pivot dari "analytics tool" ke "infrastructure protocol"
- 3 layers: On-chain Reputation (EAS) + Smart Follow + AI Agent (AgentKit)
- Composable: protocol lain (Credifi, Agently, dll) bisa query Miora scores
- Deep Base integration: EAS + AgentKit + x402

### Step 8: Competitive Analysis (Wallet Reputation Space)
- Cred Protocol — credit/lending risk (bukan trading quality)
- ChainAware — fraud detection (bukan trading performance)
- Nomis — general reputation (tidak spesifik trading)
- zScore — academic/research (belum production)
- **Temuan**: Tidak ada yang fokus di "trading quality reputation" di Base. Miora own niche ini.

### Step 9: Application Draft
- Analisis semua pertanyaan form Base Batches
- 15 pertanyaan bisa dijawab, 16 perlu diisi manual (personal info)
- 7 jawaban yang perlu di-craft hati-hati (Base-specific usage, unique insight, dll)

---

## 2. Key Findings

### Kenapa Pivot ke V2
| V1 Problem | V2 Solution |
|---|---|
| Analytics tool sudah saturated | Infrastructure protocol (new primitive) |
| Bukan Solana-native, bukan Base-native | Deep Base integration (EAS + AgentKit + x402) |
| AI wrapper bukan differentiator | Scoring engine = proprietary moat, AI = narration layer |
| No on-chain presence | EAS attestation on-chain |
| No composability | Other protocols can query scores |
| "Tool" category | "Protocol + Agent" category |

### Kenapa Base, Bukan Solana
| Factor | Solana Frontier | Base Student Track |
|---|---|---|
| Chain alignment | 5/10 — multi-chain, not Solana-native | 8.5/10 — 5/6 chains are EVM, Base included |
| Narrative fit | 5.5/10 — "beginner-friendly" bukan Solana vibe | 8/10 — matches Base "bring next billion onchain" |
| Competition | Very high — 1,500+ submissions per hackathon | More favorable — student track separate pool |
| Format | Hackathon — working demo critical | Application-based — "doesn't matter if code existed before" |

### BB003 Gap Analysis
| Category | Covered by BB003 | Gap for Miora |
|---|---|---|
| AI agent infra | Blockrun, Agently, 4Mica, Floe Labs | — |
| Lending/Credit | Credifi, Tomorrow | — |
| Privacy | OPAL | — |
| Prediction markets | Onsight, JPEG App | — |
| Neobank | Liminal | — |
| FX/Payments | Nivo | — |
| **Trading intelligence** | **Nobody** | **✅ Miora fills this** |
| **Wallet reputation** | **Nobody** | **✅ Miora fills this** |

---

## 3. External References

### Base Ecosystem
- [Base Batches Student Track — Devfolio](https://base-batches-student-track-3.devfolio.co/)
- [Introducing Base Batches 003](https://blog.base.org/introducing-base-batches-003-2)
- [Base 2026 Mission, Vision, and Strategy](https://blog.base.org/2026-mission-vision-and-strategy)
- [Base 2026 Strategy — The Defiant](https://thedefiant.io/news/blockchains/base-doubles-down-on-global-markets-stablecoins-and-ai-agents)
- [Base 2026 Strategy — CoinTelegraph](https://cointelegraph.com/news/base-joins-ethereum-tron-others-betting-big-ai-agent-future)

### Coinbase Infrastructure
- [Coinbase AgentKit Docs](https://docs.cdp.coinbase.com/agentkit/docs/add-agent-capabilities)
- [Coinbase Agentic Wallets](https://docs.cdp.coinbase.com/agentic-wallet/welcome)
- [Coinbase Agentic Wallets — CoinTelegraph](https://cointelegraph.com/news/coinbase-launches-crypto-wallets-built-ai-agents)
- [Coinbase Agentic Wallets — Bitrue](https://bitrue.com/blog/what-is-agentic-wallets-coinbase-ai-crypto-autonomy)
- [Coinbase Verifications (EAS)](https://docs.cdp.coinbase.com/verifications/introduction/welcome)
- [Coinbase Onchain Verification](https://help.coinbase.com/coinbase/getting-started/verify-my-account/onchain-verification)

### EAS (Ethereum Attestation Service)
- [EAS Official](https://attest.org/)
- [EAS — 4pillars Explainer](https://4pillars.io/en/articles/eas-the-base-layer-for-attestations)
- [EAS — QuickNode Guide](https://www.quicknode.com/guides/ethereum-development/smart-contracts/what-is-ethereum-attestation-service-and-how-to-use-it)
- [EAS SDK — npm](https://www.npmjs.com/package/@ethereum-attestation-service/eas-sdk)
- [Compliance Gating with Coinbase Verifications (EAS)](https://insumermodel.com/blog/compliance-gating-coinbase-verifications.html)

### x402 Protocol
- [KuCoin: What is x402?](https://www.kucoin.com/blog/en-what-is-x402-why-this-protocol-is-the-disruptive-backbone-for-ai-agents)
- [Turnkey: Agentic Stablecoin Micropayments](https://www.turnkey.com/blog/agentic-stablecoin-micropayments-machine-payment-protocol-x402)
- [Stellar: x402 on Stellar](https://www.stellar.org/blog/foundation-news/x402-on-stellar)
- [CoinDesk: x402 Demand Not There Yet](https://www.coindesk.com/markets/2026/03/11/coinbase-backed-ai-payments-protocol-wants-to-fix-micropayment-but-demand-is-just-not-there-yet)
- [Cred Protocol: x402 Payments for Credit API](https://credprotocol.com/blog/x402-payments-credit-api)

### Wallet Reputation / Credit Score
- [Wallet Reputation Is the New Credit Score](https://rnwy.com/blog/wallet-reputation-credit-score)
- [Web3 Reputation Score Comparison 2026](https://chainaware.ai/blog/web3-reputation-score-comparison-2026/)
- [ChainAware Wallet Rank Guide](https://chainaware.ai/blog/chainaware-wallet-rank-guide/)
- [Cred Protocol — MCP Services](https://credprotocol.com/blog/mcp-services-ai-agents-reputation-intelligence)
- [Galaxy: The New Age of Onchain Credit](https://www.galaxy.com/insights/perspectives/the-new-age-in-onchain-credit-markets)
- [zScore: Universal Decentralised Reputation System](https://arxiv.org/html/2503.05718v1)
- [On-Chain Identity as Credit Score](https://www.snowball.money/blog/the-reputation-economy-how-on-chain-identity-becomes-the-new-credit-score)

### AI Agent Economy
- [MEXC: Solana and the Agent Economy](https://blog.mexc.com/news/solana-and-the-agent-economy-will-the-future-internet-be-driven-by-ai-transactions/)
- [MEXC: How AI Agents Execute On-Chain Trades 2026](https://blog.mexc.com/autonomous-wealth-in-crypto-how-ai-agents-execute-on-chain-trades-optimize-defi-yields-and-outperform-human-traders-in-2026/)
- [Coinbase Launches Agentic Wallets — KuCoin](https://www.kucoin.com/news/flash/coinbase-launches-agentic-wallets-for-autonomous-ai-agent-commerce)
- [BingX: Top 10 Base AI Agent Projects 2026](https://bingx.com/en/learn/article/top-ai-agent-projects-in-base-ecosystem)

### BB003 Cohort Projects
- [Agently — Routing Layer for AI Agents](https://agently.to/)
- [Agently Docs — use-agently CLI](https://www.mintlify.com/AgentlyHQ/use-agently/introduction)
- [OPAL — Privacy Perp DEX](https://opaldex.com/)
- [OPAL Whitepaper](https://docs.opaldex.com/executive-summary)
- [JPEG App — Opinion Markets](https://jpeg.fun)
- [Credifi — Uncollateralized Lending](https://credi.fi)
- [Liminal — Self-custodial Neobank](https://becomeliminal.com)
- [Blockrun.ai — AI Agent Infra](https://lobehub.com/skills/ngxtm-devkit-blockrun)

### Colosseum (untuk perbandingan Solana Frontier)
- [How to Win a Colosseum Hackathon](https://blog.colosseum.org/how-to-win-a-colosseum-hackathon/)
- [Perfecting Your Hackathon Submission](https://blog.colosseum.org/perfecting-your-hackathon-submission/)
- [Announcing the Solana Frontier Hackathon](https://blog.colosseum.com/announcing-the-solana-frontier-hackathon/)
- [Colosseum Frontier Website](https://colosseum.com/frontier)

### DeFi Security (untuk konteks)
- [Solana DeFi Structural Challenges 2026](https://defi-planet.com/2026/04/solana-defi-faces-structural-challenges-as-new-yield-framework-targets-risk-transparency/)
- [Forbes: $285M Drift Hack](https://www.forbes.com/sites/jemmagreen/2026/04/11/285m-hack-proved-defis-decentralisation-promise-is-still-a-fiction/)
- [Solana Infrastructure Crisis — Validator Exodus](https://www.ainvest.com/news/solana-infrastructure-crisis-validator-exodus-slow-update-adoption-threat-network-security-long-term-2601/)
