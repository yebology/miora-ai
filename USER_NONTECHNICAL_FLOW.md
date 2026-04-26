# Miora AI — User Flow (Non-Technical)

## Apa itu Miora?

Miora membantu kamu menemukan trader terbaik di Base, lalu otomatis ikut trade mereka — dengan aturan dan budget yang kamu tentukan sendiri.

---

## Step 1: Connect Wallet

Buka Miora → klik "Connect Wallet" → pilih MetaMask → selesai. Wallet address kamu = identitas kamu di Miora. Tidak perlu email, tidak perlu password.

---

## Step 2: Analyze Wallet

Kamu lihat tweet: "Wallet 0xABC baru profit 5x dari token BRETT!"

1. Buka halaman **Analyze**
2. Paste address 0xABC
3. Klik **Analyze**
4. Tunggu ~5 detik

Miora tampilkan:
- **Score: 73/100** — "Conditional Follow"
- **Win rate: 68%** — 7 dari 10 trade profit
- **AI insight:** "Wallet ini disiplin tapi 15% trade-nya di token berisiko. Follow dengan filter."
- **Conditions:** "Hanya perhatikan token dengan liquidity > $120k dan market cap > $450k"
- **Verified on Base** — score ini di-publish on-chain, siapapun bisa verifikasi

Sekarang kamu tahu: wallet ini **oke tapi perlu hati-hati**.

---

## Step 3: Follow Wallet

Kamu mau dapat notifikasi kalau wallet ini trade:

1. Pilih conditions yang mau aktif (liquidity, market cap, dll)
2. Klik **Follow Wallet**

Sekarang kamu ada di **Watchlist**. Setiap kali wallet 0xABC trade token yang match conditions kamu, kamu dapat notifikasi real-time di browser.

---

## Step 4: Dapat Notifikasi

Jam 3 pagi, wallet 0xABC beli token PEPE.

Kamu bangun pagi, buka Miora → tab **Notifications**:
- "Wallet 0xABC bought 500M PEPE"
- AI assessment: "⚠️ PEPE is a meme token. Moderate liquidity. Only buy what you can afford to lose."

Kamu baca assessment → decide sendiri mau ikut atau tidak.

---

## Step 5: Buat Wallet Bot (Otomatis)

Kamu capek bangun pagi dan telat terus. Mau otomatis:

1. Buka halaman **Agent**
2. Klik **Wallet Bot**
3. Pilih wallet 0xABC dari watchlist
4. Set budget: 500 USDT, max per trade: 50 USDT
5. Conditions auto-filled dari hasil analyze
6. Klik **Create Bot** → MetaMask popup → deposit 500 USDT ke bot wallet → bot auto-start

Sekarang bot jalan 24/7:
- Wallet 0xABC beli LINK (liquidity $95M) → bot ikut beli 50 USDT ✅
- Wallet 0xABC beli SCAMTOKEN (liquidity $5k) → bot skip ❌ (di bawah condition)
- Wallet 0xABC jual LINK → bot ikut jual ✅

Kamu bisa lihat semua trade di **Trade History**. Pause atau stop kapan saja.

---

## Step 6: Buat Consensus Bot (Premium)

Kamu mau lebih yakin sebelum trade — bukan cuma ikut 1 wallet, tapi tunggu sampai banyak wallet bagus beli token yang sama:

1. Buka halaman **Agent**
2. Klik **Consensus Bot** (Premium)
3. Set: budget 300 USDT, max per trade 30 USDT
4. Min wallet score: 75
5. Threshold: 3 wallet harus beli token yang sama
6. Time window: 1 jam
7. Klik **Create Bot** → lalu **Start**

Bot scan semua wallet yang pernah di-analyze di Miora:
- Wallet A (score 85) beli PEPE → catat
- Wallet B (score 90) beli PEPE → catat
- Wallet C (score 78) beli PEPE → **3 wallet! Consensus!** → bot beli 30 USDT PEPE ✅

Lebih aman karena bukan cuma 1 orang yang beli — 3 trader bagus setuju.

---

## Ringkasan

| Fitur | Apa yang dilakukan | Siapa yang cocok |
|-------|-------------------|-----------------|
| Analyze | Cek apakah wallet bagus | Semua orang |
| Follow + Notifikasi | Dapat alert kalau wallet trade | Yang mau monitor manual |
| Wallet Bot | Auto-trade ikut 1 wallet | Yang mau otomatis |
| Consensus Bot | Auto-trade kalau banyak wallet setuju | Yang mau lebih aman |
