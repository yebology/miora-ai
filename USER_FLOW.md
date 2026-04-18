# Miora AI — User Flow (End to End)

## Flow 1: Analyze Wallet

```
User buka /analyze
    │
    ├── Paste wallet address (0x...)
    ├── Chain: Base (otomatis)
    ├── Pilih jumlah transaksi: 10
    ├── Klik "Analyze"
    │
    ▼
Frontend → POST /api/wallets/analyze { address, chain: "base", limit: 10 }
    │
    ▼
Backend (services/wallet.go → AnalyzeWallet):
    │
    ├── 1. Fetch transaksi dari Base via Alchemy (clients/evm.go)
    │      → Dapat 10 transaksi terakhir wallet
    │
    ├── 2. Fetch token data dari DexScreener (clients/dexscreener.go)
    │      → Liquidity, market cap, pair age, 24h volume per token
    │
    ├── 3. Fetch historical prices dari Moralis (clients/moralis.go)
    │      → Harga token saat wallet beli (by block number)
    │
    ├── 4. FIFO PnL matching (services/wallet_helper.go)
    │      → Match buy → sell per token, hitung realized + unrealized PnL
    │
    ├── 5. Scoring engine (services/scoring.go)
    │      → 6 faktor: win rate, profit consistency, entry timing,
    │        token quality, trade discipline, risk exposure
    │      → Final score 0-100
    │      → Recommendation: full_follow / conditional_follow / avoid
    │
    ├── 6. Save ke database
    │      → Wallet, Transactions, WalletMetric
    │
    ├── 7. Generate conditions (kalau conditional_follow)
    │      → Dynamic thresholds dari data wallet sendiri
    │      → "Liquidity > $120k, MCap > $450k, Pair age > 8h"
    │
    ├── 8. AI insight via Gemini (services/ai.go)
    │      → "This wallet shows disciplined trading with 68% win rate..."
    │
    └── 9. [ASYNC] Publish EAS attestation (clients/eas.go)
           → ABI-encode: (score, recommendation, totalTxns, chain)
           → Sign + send tx ke EAS contract di Base Sepolia
           → Dapat attestation UID + tx hash
           → Update WalletMetric di database
    │
    ▼
Frontend tampilkan:
    ├── Score ring (73/100)
    ├── Recommendation badge (Conditional Follow)
    ├── 6 metric bars (win rate, profit consistency, dll)
    ├── AI insight card (plain language explanation)
    ├── Conditions card (checkboxes untuk follow)
    ├── Traded tokens table (25 tokens dengan PnL)
    └── Attestation badge → link ke BaseScan EAS explorer
```

---

## Flow 2: Follow Wallet (Smart Follow)

```
User di halaman analyze result
    │
    ├── Lihat recommendation: "Conditional Follow"
    ├── Pilih conditions: ✅ Min liquidity, ✅ Min mcap
    ├── Klik "Follow Wallet"
    │
    ▼
Frontend cek: user sudah login?
    ├── Belum → tampilkan Auth Guard Modal → "Sign in with Google"
    └── Sudah → lanjut
    │
    ▼
Frontend → POST /api/watchlist/follow {
    wallet_address, chain: "base",
    recommendation: "conditional_follow",
    conditions: ["min_liquidity", "min_mcap"],
    email_notify: true
}
    │
    ▼
Backend (handlers/watchlist.go → Follow):
    ├── Verify Firebase token
    ├── Check wallet belum di-follow
    ├── Save ke database (entities/watchlist.go)
    └── Return success
    │
    ▼
User sekarang bisa lihat wallet ini di /watchlist
```

---

## Flow 3: Real-time Trade Notification

```
Background (services/monitor.go → Start):
    │
    ├── Polling setiap 30 detik
    ├── Untuk setiap wallet yang di-follow:
    │
    ├── 1. Fetch transaksi terbaru via Alchemy
    │      → Bandingkan dengan count sebelumnya
    │      → Ada transaksi baru? Lanjut.
    │
    ├── 2. Fetch token data dari DexScreener
    │      → Liquidity, market cap, pair age
    │
    ├── 3. Check conditions (monitor_helper.go)
    │      → User pilih min_liquidity: token liquidity > $100k? ✅
    │      → User pilih min_mcap: token mcap > $500k? ✅
    │      → Semua pass → notify
    │
    ├── 4. AI risk assessment via Gemini
    │      → "⚠️ PEPE is a meme token with moderate liquidity..."
    │
    ├── 5. Send WebSocket notification (ws/hub.go)
    │      → Real-time push ke browser user
    │
    ├── 6. Save notification ke database
    │      → Notification entity (token, direction, value, AI assessment)
    │
    └── 7. Send email via Resend (clients/resend.go)
           → Async, non-blocking
    │
    ▼
User di /watchlist tab "Notifications":
    ├── Lihat: "Wallet 0x1234... bought 500M PEPE"
    ├── AI assessment: "⚠️ Moderate liquidity, high risk"
    └── Bisa klik untuk lihat detail wallet
```

---

## Flow 4: AI Trading Agent

```
User buka /agent
    │
    ├── Login required → Sign in with Google
    │
    ▼
Agent Status Card:
    ├── Status: Paused
    ├── Budget: $0 (belum diset)
    ├── Trades: 0
    │
    ▼
User isi Agent Config Form:
    ├── Budget: $500
    ├── Max per trade: $50
    ├── Risk tolerance: Medium
    ├── Min score: 70
    ├── Conditions: ✅ Min liquidity, ✅ Min mcap
    ├── Klik "Save Configuration"
    │
    ▼
Frontend → PUT /api/agent/config { budget: 500, max_per_trade: 50, ... }
    │
    ▼
User klik "Start Agent"
    │
    ▼
Frontend → POST /api/agent/start
Backend set status = "active"
    │
    ▼
Background (services/agent_loop.go → Start):
    │
    ├── Polling setiap 30 detik
    ├── Get semua agent configs yang status = "active"
    │
    ├── Untuk setiap active agent:
    │   │
    │   ├── 1. Cari wallet dengan score >= 70 (minScore)
    │   │
    │   ├── 2. Cek transaksi baru dari wallet tersebut
    │   │      → "Wallet 0xABC (score 87) baru beli PEPE"
    │   │
    │   ├── 3. Evaluate conditions
    │   │      → Liquidity > $100k? ✅
    │   │      → MCap > $500k? ✅
    │   │
    │   ├── 4. Check budget
    │   │      → $500 - $120 spent = $380 remaining
    │   │      → $50 per trade < $380? ✅
    │   │
    │   ├── 5. AI risk assessment via Gemini
    │   │      → "Low risk — high liquidity, established token"
    │   │
    │   ├── 6. Risk tolerance check
    │   │      → Medium tolerance + liquidity > $100k? ✅
    │   │
    │   ├── 7. Execute swap via Python sidecar
    │   │      Go backend → POST http://localhost:8090/swap
    │   │      Python sidecar → AgentKit → Agentic Wallet → on-chain tx
    │   │      → Return tx hash
    │   │
    │   ├── 8. Record trade di database
    │   │      → AgentTrade { status: "executed", token: "PEPE", amount: $50 }
    │   │
    │   └── 9. Update budget
    │          → total_spent += $50, total_trades += 1
    │
    ▼
User lihat di /agent:
    ├── Status: Active
    ├── Remaining: $330
    ├── Trades: 5
    └── Trade History:
        ├── ✅ PEPE — $50 — "Bought because wallet 0xABC (score 87) bought it"
        ├── ⏭️ NEWTOKEN — skipped — "Liquidity $8k below min $100k"
        └── ✅ LINK — $30 — "Bought because wallet 0xABC (score 87) bought it"
```

---

## Flow 5: Query On-chain Reputation

```
Any protocol / AI agent / user
    │
    ├── GET /api/reputation/0x1234...
    │
    ▼
Backend (handlers/reputation.go → GetReputation):
    │
    ├── 1. Lookup wallet di database
    ├── 2. Get WalletMetric (score, attestation UID)
    ├── 3. Query EAS contract on-chain (clients/eas.go → GetAttestation)
    │      → Decode attestation data: score, recommendation, totalTxns, chain
    │
    ▼
Response:
{
    "address": "0x1234...",
    "score": 73,
    "recommendation": "conditional_follow",
    "total_transactions": 47,
    "attestation_uid": "0xabc...",
    "attester": "0xMioraWallet...",
    "timestamp": 1713200000,
    "explorer_url": "https://base-sepolia.easscan.org/attestation/view/0xabc..."
}
    │
    ▼
Protocol bisa pakai score ini untuk:
    ├── Lending: "Score > 70 → lower interest rate"
    ├── Airdrop: "Score > 80 → eligible"
    ├── AI agent: "Score > 60 → safe to copy trade"
    └── DEX: "Score < 30 → show warning"
```

---

## Flow 6: Sign In

```
User klik "Sign In" di navbar
    │
    ▼
Auth Guard Modal muncul
    ├── Klik "Sign in with Google"
    │
    ▼
Firebase Auth → Google OAuth popup
    ├── User pilih Google account
    ├── Firebase return ID token
    │
    ▼
Frontend → GET /api/auth/me (with Firebase token)
    │
    ▼
Backend (handlers/auth.go → Me):
    ├── Verify Firebase token (middleware/firebase.go)
    ├── Extract: firebase_uid, email, name
    ├── Find or create user di database
    └── Return user data
    │
    ▼
Frontend simpan user di AuthContext
    ├── Navbar: "Sign In" → avatar + name
    ├── Protected pages (watchlist, agent) sekarang accessible
    └── Firebase token disimpan untuk API calls
```
