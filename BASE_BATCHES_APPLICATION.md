# Base Batches 003: Student Track — Application Draft

---

## Company Name

Miora AI

---

## Website / Product URL

https://github.com/yebology/miora-ai

---

## Demo URL

https://drive.google.com/file/d/1V73h7VNh0EZggx8mVDi0BXFlDTV8oL_1/view?usp=sharing

---

## Describe what your company does (~50 chars)

Trading reputation protocol + AI bots on Base

---

## What is your product's unique value proposition?

Three things that no other tool in the market does:

1. On-chain trading reputation — Every wallet score is published as an EAS attestation on Base. This is not just a number in our database. It is a composable on-chain primitive that any protocol, agent, or dApp can query and build on top of. No other wallet analytics tool publishes scores on-chain.

2. Consensus Bot — Most copy trading tools let you follow one wallet. Miora's Consensus Bot watches all analyzed wallets and only trades when multiple high-score wallets buy the same token within a short time window. This crowd intelligence approach gives much higher confidence than following a single trader.

3. Dynamic guard rails from the wallet's own data — When Miora recommends "Conditional Follow," the filter conditions (minimum liquidity, market cap, volume, pair age) are computed from that specific wallet's trading history, not hardcoded thresholds. Every wallet gets personalized guard rails based on its own behavior.

---

## What part of your product is onchain?

EAS attestation (trading reputation scores) and AI trading bots (Coinbase AgentKit + CDP Server Wallet) on Base Sepolia.

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

Surabaya, East Java, Indonesia

---

## Do you already have a token?

No.

---

## What part of your product uses Base?

Everything is exclusive to Base. EAS attestation for on-chain reputation, AI trading bots via Coinbase AgentKit, and wallet analysis all run on Base Sepolia. No other network is supported.

---

## Founder(s) Names and Contact Information

Yobel Nathaniel Filipus || yobelnathaniel12@gmail.com

---

## Please describe each founder's background and add their LinkedIn profile(s)

Yobel Nathaniel Filipus — Final year undergraduate student at Universitas Ciputra Surabaya. Solo founder handling all technical development (backend, frontend, agent sidecar, infrastructure).

LinkedIn: https://linkedin.com/in/yobelnathanielfilipus

---
## Please enter the URL of a ~1-minute unlisted video introducing the founder(s) and what you're building

https://drive.google.com/file/d/1NrBd72H5XUzvQeBCYdTPXT7oHo86flpK/view?usp=sharing

---

## Who writes code or handles technical development?

I (Yobel) handle all technical development — backend (Go + Fiber), frontend (Next.js + TypeScript), agent sidecar (Python + AgentKit), smart contracts (Foundry), and infrastructure (Docker).

---

## How long have the founders known each other and how did you meet?

Solo founder.

---

## How far along are you?

Prototype

---

## How long have you been working on this?

3 weeks.

---

## How much of that time is full-time vs part-time?

Part-time. Building alongside final year university coursework.

---

## What part of your product is magic or impressive?

The scoring engine analyzes any wallet using FIFO PnL matching and 5 metrics, then gives a clear recommendation: Full Follow, Conditional Follow with personalized guard rails, or Avoid. If the score is good, users deploy an AI bot via Coinbase AgentKit that copies trades automatically. The Consensus Bot only trades when multiple high-score wallets buy the same token, giving higher confidence than following one trader.

---

## What is your unique insight or advantage?

Existing tools show wallet data. Miora shows a decision (Follow, Conditional Follow, or Avoid) and publishes it on-chain via EAS so anyone can verify. Then users can deploy AI bots that act on that decision automatically. No other tool combines on-chain verifiable reputation with autonomous copy-trading bots.

---

## Do you plan on raising capital from VCs? Do you plan to launch a token?

Open to raising capital to go full-time and scale to Base mainnet. Token launch is planned but currently focused on improving the platform first.

---

## Do you have users or customers?

Not yet — currently in prototype stage. Plan to onboard beta testers on Base before public launch.

---

## Active users / paying customers

N/A — pre-launch.

---

## Revenue

N/A — pre-launch. Planned revenue model:
- Consensus Bot as premium feature (subscription or per-trade fee) — trades with higher confidence using multi-wallet agreement signals
- Reputation API at scale — public endpoint `GET /api/reputation/:address` is free now, monetizable via subscription or per-query fees when demand grows

---

## Dune dashboards / deployed contract addresses

- EAS Schema UID: `0xe8287a3e4882e5f061ea2fb9e7b0810ce94740005aaa46497fcf1712ff5a41b0`
- EAS Schema Explorer: https://base-sepolia.easscan.org/schema/view/0xe8287a3e4882e5f061ea2fb9e7b0810ce94740005aaa46497fcf1712ff5a41b0

---

## Why do you want to join Base Batches?

Three reasons:

1. Mentorship and feedback — As a solo undergraduate founder, direct access to experienced builders and the Base team would accelerate product decisions that currently take me weeks of research. Specifically, I want guidance on scaling the scoring engine and optimizing bot execution for mainnet.
2. Base ecosystem alignment — Miora is built exclusively on Base's own infrastructure (EAS + AgentKit). Being embedded in the Base ecosystem means better integration support, early access to new tools, and distribution to Base's 34M+ monthly active users.
3. Demo Day exposure — Presenting to investors and the Base community in San Francisco would be a transformative opportunity to validate the product direction and explore funding for full-time development after graduation.

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

- [x] Record 1-minute founder video
- [x] Register EAS schema on Base Sepolia (`make register-schema`)
- [ ] Fill all [FILL] fields
- [x] Paste EAS schema UID + MockUSDT contract in Dune/contract addresses field
- [x] Make GitHub repo public
- [ ] Review all answers for clarity and conciseness
- [ ] Submit before April 27, 2026
