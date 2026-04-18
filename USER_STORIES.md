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
5. Miora minta login → Andi sign in with Google
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

## Story 4: Andi setup AI trading agent

> "Sebagai Andi, saya ingin agent yang trade otomatis berdasarkan wallet terbaik, supaya saya tidak perlu monitor 24/7."

**Scenario:**

1. Andi buka halaman Agent
2. Isi konfigurasi:
   - Budget: $200
   - Max per trade: $20
   - Risk tolerance: Medium
   - Min score: 75 (hanya follow wallet score tinggi)
   - Conditions: ✅ Min liquidity, ✅ Min mcap
3. Klik "Save Configuration"
4. Klik "Start Agent"
5. Agent mulai jalan di background
6. 2 jam kemudian:
   - Wallet 0xDEF (score 88) beli LINK
   - Agent evaluate: liquidity $95M ✅, mcap $8.5B ✅, budget cukup ✅
   - AI assessment: "Very low risk — top-tier token"
   - Agent execute: beli $20 LINK via Agentic Wallet
7. Andi buka Agent page → lihat trade history:
   - "✅ LINK — $20 — Bought because wallet 0xDEF (score 88) bought it"
8. Keesokan harinya:
   - Wallet 0xGHI (score 65) beli NEWTOKEN
   - Agent evaluate: score 65 < minScore 75 ❌
   - Agent skip: "Wallet score below threshold"
   - Trade dicatat sebagai "skipped"

**Acceptance criteria:**
- Agent hanya trade kalau SEMUA conditions pass
- Budget tidak pernah terlampaui
- Setiap trade (executed/skipped/failed) dicatat dengan alasan
- User bisa pause/resume kapan saja

---

## Story 5: Andi cek reputasi on-chain

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

## Story 6: Protocol lain query Miora score

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

## Story 7: Andi pause agent karena market crash

> "Sebagai Andi, saya ingin bisa pause agent kapan saja kalau market sedang tidak stabil, supaya agent tidak trade di kondisi buruk."

**Scenario:**

1. Andi lihat berita: "Base DeFi market crash 20%"
2. Buka halaman Agent
3. Status: Active, 8 trades executed, $160 spent
4. Klik "Pause Agent"
5. Agent langsung berhenti polling
6. Tidak ada trade baru yang dieksekusi
7. 3 hari kemudian, market stabil
8. Andi buka Agent → klik "Start Agent"
9. Agent resume polling dan trading

**Acceptance criteria:**
- Pause immediate — tidak ada trade setelah pause
- State preserved — budget, trade history tetap ada
- Resume tanpa perlu reconfigure
