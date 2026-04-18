# Miora AI — User Stories

## Persona

**Andi** — Mahasiswa 22 tahun, baru mulai trading di Base 3 bulan lalu. Punya $500 di wallet MetaMask. Sering lihat "whale alerts" di Twitter tapi tidak tahu cara evaluasi apakah wallet itu benar-benar bagus. Pernah rugi $150 karena ikut-ikutan beli token yang ternyata rug pull.

---

## Story 1: Andi menemukan trader bagus

> "Sebagai Andi, saya ingin tahu apakah wallet yang sering disebut di Twitter benar-benar trader yang bagus, supaya saya tidak asal ikut-ikutan."

**Scenario:**

1. Andi lihat tweet: "Wallet 0xABC just made 5x on $BRETT 🔥"
2. Andi buka Miora → halaman Analyze
3. Paste address 0xABC, klik Analyze
4. Tunggu ~5 detik
5. Miora tampilkan:
   - Score: **73/100** — Conditional Follow
   - Win rate: 68% — "Lumayan, tapi bukan yang terbaik"
   - Risk exposure: 15% — "Ada beberapa trade di token low-liquidity"
   - AI insight: "Wallet ini disiplin dengan 68% win rate. Tapi 15% trade-nya di token berisiko. Follow dengan filter — hanya perhatikan token dengan liquidity di atas $100k."
6. Andi sekarang tahu: wallet ini **oke tapi perlu hati-hati**, bukan blindly follow

**Acceptance criteria:**
- Score 0-100 ditampilkan dalam < 10 detik
- AI insight dalam bahasa yang mudah dipahami
- Recommendation jelas: Full Follow / Conditional Follow / Avoid

---

## Story 2: Andi follow wallet dengan conditions

> "Sebagai Andi, saya ingin follow wallet yang bagus tapi hanya dapat notifikasi untuk trade yang aman, supaya saya tidak tergoda beli token sampah."

**Scenario:**

1. Andi lihat hasil analyze: Conditional Follow (score 73)
2. Miora tampilkan 4 conditions:
   - ✅ Token liquidity > $120k
   - ✅ Market cap > $450k
   - ☐ Pair age > 8 hours
   - ☐ 24h volume > $60k
3. Andi centang liquidity dan market cap
4. Klik "Follow Wallet"
5. Miora minta connect wallet → Andi connect MetaMask
6. Wallet ditambahkan ke watchlist
7. Andi sekarang hanya dapat notifikasi kalau wallet 0xABC beli token yang liquidity > $120k DAN mcap > $450k

**Acceptance criteria:**
- Conditions computed dari data wallet sendiri (bukan hardcoded)
- User bisa pilih conditions mana yang aktif
- Hanya trade yang match conditions yang trigger notifikasi

---

## Story 3: Andi dapat notifikasi real-time

> "Sebagai Andi, saya ingin langsung tahu kalau wallet yang saya follow melakukan trade, supaya saya bisa ambil keputusan cepat."

**Scenario:**

1. Andi sedang browsing di Miora (tab watchlist terbuka)
2. Wallet 0xABC beli 500M PEPE di Base
3. Miora detect transaksi baru (polling setiap 30 detik)
4. Miora cek conditions:
   - PEPE liquidity: $250k > $120k ✅
   - PEPE mcap: $4.2M > $450k ✅
5. Miora generate AI risk assessment:
   - "⚠️ PEPE is a meme token with moderate liquidity. It can swing wildly — only buy what you can afford to lose."
6. Notifikasi muncul di browser Andi (WebSocket):
   - "Wallet 0xABC bought 500M PEPE"
   - AI assessment ditampilkan
7. Andi juga dapat email notification
8. Andi baca AI assessment → decide untuk **tidak ikut beli** karena meme token

**Acceptance criteria:**
- Notifikasi muncul < 1 menit setelah trade terdeteksi
- AI risk assessment ada di setiap notifikasi
- Email notification dikirim kalau email_notify aktif

---

## Story 4: Andi buat wallet bot dari watchlist

> "Sebagai Andi, saya ingin bot yang otomatis copy trade wallet terbaik dari watchlist saya, supaya saya tidak perlu monitor 24/7."

**Scenario:**

1. Andi buka halaman Agent (Bot Management)
2. Klik "Create Bot" → pilih type: **Wallet Bot**
3. Miora tampilkan daftar wallet dari watchlist Andi
4. Andi pilih wallet 0xDEF (score 88, Full Follow)
5. Conditions otomatis terisi dari hasil analyze wallet 0xDEF:
   - ✅ Min liquidity: $100k
   - ✅ Min mcap: $1M
   - ✅ Min pair age: 6 hours
6. Andi isi parameter:
   - Budget: $200
   - Max per trade: $20
   - Min score: 75
7. Klik "Create Bot" → bot dibuat
8. Klik "Start" → bot mulai jalan di background
9. 2 jam kemudian:
   - Wallet 0xDEF beli LINK
   - Bot evaluate: liquidity $95M ✅, mcap $8.5B ✅, pair age OK ✅, budget cukup ✅
   - AI assessment: "Very low risk — top-tier token"
   - Bot execute: beli $20 LINK via Agentic Wallet
10. Andi buka Bot page → lihat trade history:
    - "✅ LINK — $20 — Bought because wallet 0xDEF (score 88) bought it"
11. 3 jam kemudian:
    - Wallet 0xDEF jual LINK
    - Bot evaluate: conditions pass ✅
    - Bot execute: jual LINK via Agentic Wallet
    - "✅ LINK — Sold because wallet 0xDEF sold it"

**Acceptance criteria:**
- Bot hanya trade kalau SEMUA conditions pass
- Bot follows both buys AND sells of target wallet
- Budget tidak pernah terlampaui
- Setiap trade (executed/skipped/failed) dicatat dengan alasan
- User bisa pause/resume kapan saja
- Conditions auto-filled dari analyze result saat create bot dari watchlist

---

## Story 5: Andi buat consensus bot

> "Sebagai Andi, saya ingin bot yang trade otomatis kalau banyak wallet bagus beli token yang sama, supaya saya dapat sinyal yang lebih kuat."

**Scenario:**

1. Andi buka halaman Agent → klik "Create Bot" → pilih type: **Consensus Bot**
2. Isi parameter:
   - Budget: $300
   - Max per trade: $30
   - Min score: 70 (hanya pertimbangkan wallet score tinggi)
   - Consensus threshold: 3 (minimal 3 wallet harus beli token yang sama)
   - Time window: 60 menit (dalam 60 menit terakhir)
3. Klik "Create Bot" → bot dibuat
4. Klik "Start" → bot mulai scan semua wallet Miora
5. 4 jam kemudian:
   - Wallet 0xAAA (score 85) beli TOKEN_X
   - Wallet 0xBBB (score 78) beli TOKEN_X
   - Wallet 0xCCC (score 91) beli TOKEN_X
   - 3 wallet dalam 60 menit → threshold tercapai ✅
   - Bot evaluate conditions: liquidity OK ✅, mcap OK ✅
   - AI assessment: "Multiple high-score wallets buying — strong consensus signal"
   - Bot execute: beli $30 TOKEN_X
6. Andi lihat di trade history:
   - "✅ TOKEN_X — $30 — Consensus: 3 wallets (0xAAA, 0xBBB, 0xCCC) bought within 60min"

**Acceptance criteria:**
- Consensus bot is a separate bot type (not a toggle inside wallet bot)
- Bot only trades when threshold number of wallets buy same token within time window
- Each wallet must meet min_score requirement
- Budget dan max per trade dipatuhi
- Trade history menunjukkan wallet mana yang trigger consensus

---

## Story 6: Andi cek reputasi on-chain

> "Sebagai Andi, saya ingin buktikan ke teman bahwa wallet saya punya trading score bagus, dengan bukti on-chain yang tidak bisa dipalsukan."

**Scenario:**

1. Andi analyze wallet sendiri di Miora
2. Score: 82 — Full Follow
3. Miora publish attestation ke Base Sepolia via EAS
4. Andi lihat badge "Verified on Base" di hasil analyze
5. Klik badge → buka BaseScan EAS explorer
6. Di explorer terlihat:
   - Attester: 0xMioraWallet
   - Recipient: 0xAndiWallet
   - Data: score 82, recommendation "full_follow", 35 transactions
   - Timestamp: 2026-04-15
7. Andi share link BaseScan ke teman: "Lihat, Miora score saya 82 — verified on-chain"
8. Teman bisa verify sendiri tanpa perlu percaya Andi

**Acceptance criteria:**
- Attestation visible di BaseScan EAS explorer
- Data bisa di-decode (score, recommendation, totalTxns, chain)
- Link explorer bisa di-share

---

## Story 7: Protocol lain query Miora score

> "Sebagai developer lending protocol di Base, saya ingin query trading reputation wallet sebelum approve pinjaman, supaya saya bisa assess risiko borrower."

**Scenario:**

1. Lending protocol integrate Miora API
2. User apply pinjaman di lending protocol
3. Lending protocol call: `GET /api/reputation/0xBorrower`
4. Miora return:
   ```json
   {
     "score": 85,
     "recommendation": "full_follow",
     "attestation_uid": "0xabc...",
     "explorer_url": "https://base-sepolia.easscan.org/..."
   }
   ```
5. Lending protocol logic:
   - Score > 80 → approve dengan bunga rendah (5%)
   - Score 50-80 → approve dengan bunga standar (10%)
   - Score < 50 → reject atau bunga tinggi (20%)
6. User dengan Miora score 85 dapat bunga 5%

**Acceptance criteria:**
- API response < 500ms
- Data konsisten dengan on-chain attestation
- Tidak perlu auth — public endpoint

---

## Story 8: Andi pause bot karena market crash

> "Sebagai Andi, saya ingin bisa pause bot kapan saja kalau market sedang tidak stabil, supaya bot tidak trade di kondisi buruk."

**Scenario:**

1. Andi lihat berita: "Base DeFi market crash 20%"
2. Buka halaman Agent
3. Andi punya 2 bot:
   - Wallet Bot (0xDEF): Active, 8 trades executed, $160 spent
   - Consensus Bot: Active, 3 trades executed, $90 spent
4. Klik "Pause" di kedua bot
5. Kedua bot langsung berhenti polling
6. Tidak ada trade baru yang dieksekusi
7. 3 hari kemudian, market stabil
8. Andi buka Agent → klik "Start" di kedua bot
9. Bot resume polling dan trading

**Acceptance criteria:**
- Pause immediate — tidak ada trade setelah pause
- State preserved — budget, trade history tetap ada
- Resume tanpa perlu reconfigure
- Each bot can be paused/started independently
