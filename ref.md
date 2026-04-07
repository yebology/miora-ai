# Recommendation System — Breakdown

## Full Follow (Score 80–100)

Wallet ini konsisten profit, low risk, dan punya strategi jelas.

User experience:
- Badge hijau "Full Follow"
- Tampilkan semua token yang wallet trade beserta PnL masing-masing
- Tombol "Copy Trade" per token — klik langsung ke swap quote
- AI insight menjelaskan kenapa wallet ini layak diikuti
- Opsi "Follow Wallet" — save ke watchlist, nanti bisa dapet notif kalau wallet trade lagi

Contoh UI:
```
🟢 Full Follow — Score: 85/100

"This wallet is a disciplined trader that focuses on 3-4 tokens with 
consistent 20-30% gains. Safe to follow."

Recent Trades:
┌─────────┬──────────┬────────────┐
│ Token   │ PnL      │ Action     │
├─────────┼──────────┼────────────┤
│ PEPE    │ +25%     │ [Buy PEPE] │
│ WIF     │ +18%     │ [Buy WIF]  │
│ BONK    │ +32%     │ [Buy BONK] │
└─────────┴──────────┴────────────┘

[Follow Wallet]  [View All Trades]
```

---

## Partial Follow (Score 60–79)

Wallet ini profitable tapi ada risiko. Nggak semua trade-nya bagus.

User experience:
- Badge kuning "Partial Follow"
- Tampilkan token yang profitable saja (filter PnL > 0)
- Token yang rugi ditampilkan tapi di-grey out dengan warning
- Tombol "Buy" hanya untuk token profitable
- AI insight menjelaskan mana yang bagus dan mana yang harus dihindari

Contoh UI:
```
🟡 Partial Follow — Score: 68/100

"This wallet wins 60% of trades but sometimes picks risky tokens. 
Follow the profitable ones, skip the rest."

Recent Trades:
┌─────────┬──────────┬────────────────────┐
│ Token   │ PnL      │ Action             │
├─────────┼──────────┼────────────────────┤
│ PEPE    │ +45%     │ [Buy PEPE]         │
│ WIF     │ +12%     │ [Buy WIF]          │
│ SCAM    │ -80%     │ ⚠️ High Risk       │
│ BONK    │ +8%      │ [Buy BONK]         │
└─────────┴──────────┴────────────────────┘

[Watch Wallet]  [View All Trades]
```

---

## Conditional Follow (Score 40–59)

Wallet ini mixed — kadang profit, kadang rugi besar. Ikuti hanya dengan syarat tertentu.

User experience:
- Badge oranye "Conditional Follow"
- Tampilkan semua token tapi TANPA tombol "Buy" langsung
- Tampilkan conditions dari AI: "Follow only if entry timing < 24h" atau "Only buy tokens with mcap > $1M"
- Tombol "Get Quote" (bukan "Buy") — user harus review dulu
- Warning banner di atas

Contoh UI:
```
🟠 Conditional Follow — Score: 48/100

"This wallet enters tokens early but exits too late, causing losses. 
Only follow if you can monitor closely and exit within 6 hours."

⚠️ Conditions:
- Only follow early entries (< 24h after token launch)
- Avoid tokens with liquidity below $50k
- Set stop-loss at 15%

Recent Trades:
┌─────────┬──────────┬────────────────────┐
│ Token   │ PnL      │ Action             │
├─────────┼──────────┼────────────────────┤
│ PEPE    │ +120%    │ [Get Quote]        │
│ RUG     │ -95%     │ ⚠️ Rug Pull        │
│ WIF     │ +30%     │ [Get Quote]        │
│ FAKE    │ -70%     │ ⚠️ Low Liquidity   │
└─────────┴──────────┴────────────────────┘

[Watch Wallet]
```

---

## Avoid (Score < 40)

Wallet ini high risk, poor performance, atau suspicious behavior.

User experience:
- Badge merah "Avoid"
- TIDAK ada tombol "Buy" atau "Get Quote"
- Tampilkan data sebagai informasi saja
- Warning besar di atas
- AI insight menjelaskan kenapa harus dihindari

Contoh UI:
```
🔴 Avoid — Score: 22/100

"This wallet trades mostly low-liquidity tokens and has lost money on 
80% of trades. Three of the tokens it traded have been flagged as 
potential rug pulls. Do not follow."

⛔ Warning: This wallet shows high-risk trading behavior.

Recent Trades:
┌─────────┬──────────┬────────────────────┐
│ Token   │ PnL      │ Status             │
├─────────┼──────────┼────────────────────┤
│ SCAM1   │ -90%     │ 🚩 Rug Pull        │
│ SCAM2   │ -85%     │ 🚩 Low Liquidity   │
│ PEPE    │ +5%      │ Minimal Profit     │
│ SCAM3   │ -100%    │ 🚩 Token Dead      │
└─────────┴──────────┴────────────────────┘

[Report Wallet]
```

---

## Data Flow

```
User input wallet address
        │
        ▼
POST /api/wallets/analyze
        │
        ▼
Backend returns:
- scoring (win_rate, profit_consistency, etc.)
- recommendation (full_follow / partial_follow / conditional_follow / avoid)
- ai_insight (natural language explanation)
- traded_tokens (list of tokens with PnL) ← needs to be added
        │
        ▼
Frontend reads recommendation → renders appropriate UI
        │
        ▼
User clicks "Buy [TOKEN]"
        │
        ▼
POST /api/swap/quote
        │
        ▼
Frontend shows quote → user confirms → execute swap (future)
```

---

## What Backend Needs to Add

Response dari `/api/wallets/analyze` perlu ditambahin field `traded_tokens`:

```json
{
  "traded_tokens": [
    {
      "contract_address": "0x6982508...",
      "symbol": "PEPE",
      "chain": "ethereum",
      "pnl_percent": 25.5,
      "buy_price": 0.0000012,
      "exit_price": 0.0000015,
      "status": "realized"
    }
  ]
}
```

Ini yang frontend butuhkan buat render list token + tombol Buy.
