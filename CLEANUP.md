# Miora AI — Codebase Cleanup (V2 Pivot)

Daftar semua file dan kode yang harus dihapus atau diubah karena pivot ke V2 (Base-only, no Solana).

---

## 1. Backend — Hapus File

| File | Alasan |
|---|---|
| `backend/app/clients/solana.go` | Alchemy Solana client — tidak dipakai lagi |
| `backend/app/clients/birdeye.go` | Birdeye Solana historical prices — tidak dipakai lagi |
| `backend/app/clients/jupiter.go` | Jupiter Solana swap quotes — tidak dipakai lagi |
| `backend/app/interfaces/birdeye.go` | IBirdeye interface — tidak dipakai lagi |

---

## 2. Backend — Edit File

### `backend/config/config.go`
- [ ] Hapus field `BirdeyeAPIKey` dari struct `Config`
- [ ] Hapus `"BIRDEYE_API_KEY"` dari array `required` di `LoadConfig()`
- [ ] Hapus baris `BirdeyeAPIKey: os.Getenv("BIRDEYE_API_KEY")` dari return value
- [ ] Update comment `AlchemyAPIKey` — hapus "& Solana"

### `backend/constants/chains.go`
- [ ] Hapus entry `"solana"` dari map `SupportedChains`
- [ ] Hapus case `"svm"` → `"solana"` dari `GetChainConfig()`
- [ ] Hapus function `IsSolana()` seluruhnya
- [ ] Update comment "Supports Ethereum mainnet, L2s..." — hapus "and Solana"

### `backend/constants/limits.go`
- [ ] Hapus `SVMTransactionLimits` variable seluruhnya
- [ ] Update `GetTransactionLimits()` — hapus `if IsSolana(chain)` branch, langsung return `EVMTransactionLimits`

### `backend/constants/error.go`
- [ ] Ubah `UnsupportedChain = "Chain must be 'evm' or 'svm'."` → `"Unsupported chain."`

### `backend/app/services/wallet.go`
- [ ] Hapus field `svmClient` dari struct `WalletService`
- [ ] Hapus field `birdeye` dari struct `WalletService`
- [ ] Hapus parameter `svmClient interfaces.BlockchainClient` dari `NewWalletService()`
- [ ] Hapus parameter `birdeye interfaces.IBirdeye` dari `NewWalletService()`
- [ ] Hapus assignment `svmClient: svmClient` dari return value
- [ ] Hapus assignment `birdeye: birdeye` dari return value

### `backend/app/services/wallet_helper.go`
- [ ] Di function `getClient()` — hapus `if constants.IsSolana(chain)` branch (return `s.svmClient`)
- [ ] Di function `getPrice()` — hapus `if constants.IsSolana(chain)` branch (Birdeye call), hanya sisakan Moralis

### `backend/app/services/swap.go`
- [ ] Hapus field `jupiter` dari struct `SwapService`
- [ ] Hapus parameter `jupiter` dari `NewSwapService()`
- [ ] Hapus `if constants.IsSolana(chain)` block di `GetQuote()` (Jupiter call)

### `backend/app/services/monitor.go`
- [ ] Hapus field `svmClient` dari struct `MonitorService`
- [ ] Hapus parameter `svmClient interfaces.BlockchainClient` dari `NewMonitorService()`
- [ ] Hapus assignment `svmClient: svmClient` dari return value

### `backend/app/services/monitor_helper.go`
- [ ] Di function `checkWallet()` — hapus `else if constants.IsSolana(chain)` branch

### `backend/router/container.go`
- [ ] Hapus parameter `birdeyeAPIKey` dari `NewContainer()` signature
- [ ] Hapus `svmClient := clients.NewAlchemySolana(alchemyAPIKey)`
- [ ] Hapus `birdeye := clients.NewBirdeye(birdeyeAPIKey)`
- [ ] Hapus `jupiter := clients.NewJupiter()`
- [ ] Update `NewWalletService()` call — hapus `svmClient` dan `birdeye` arguments
- [ ] Update `NewSwapService()` call — hapus `jupiter` argument (hanya `oneInch`)
- [ ] Update `NewMonitorService()` call — hapus `svmClient` argument

### `backend/router/routes.go`
- [ ] Update `NewContainer()` call — hapus `cfg.BirdeyeAPIKey` argument

### `backend/app/handlers/wallet.go`
- [ ] Update comment `"chain": "evm"` — hapus `// required — must be "evm" or "svm"`

### `backend/app/http/swap.go`
- [ ] Update comment — hapus `from Jupiter (Solana) or`

### `backend/.env`
- [ ] Hapus `BIRDEYE_API_KEY=...`

---

## 3. Frontend — Hapus File

| File | Alasan |
|---|---|
| `frontend/public/chains/solana.svg` | Solana chain logo — tidak dipakai lagi |

---

## 4. Frontend — Edit File

### `frontend/constants/landing.ts`
- [ ] Update features "DEX Aggregator" description — hapus "Jupiter (Solana) and", ubah ke "Best price routing across Base DEXs via 1inch"
- [ ] Update features "Multi-Chain Support" description — hapus "and Solana", ubah ke "Base, Ethereum, Arbitrum, Optimism, and Polygon"
- [ ] Hapus `{ name: "Solana", logo: "/chains/solana.svg", color: "#14F195" }` dari chains array

### `frontend/components/analyze/analyze-form.tsx`
- [ ] Hapus `{ value: "solana", label: "Solana" }` dari CHAINS array
- [ ] Hapus `SVM_LIMITS` variable seluruhnya
- [ ] Update `getLimits()` — langsung return `EVM_LIMITS` (tidak perlu conditional)

### `frontend/app/swap/page.tsx`
- [ ] Hapus `{ value: "solana", label: "Solana" }` dari chains array
- [ ] Hapus `chain === "solana"` conditional di route mock — langsung pakai "Uniswap V3 → SushiSwap"
- [ ] Update subtitle — hapus "Jupiter and", ubah ke "Best price routing via 1inch on Base"
- [ ] Hapus `chain === "solana" ? "Jupiter" : "1inch"` conditional — langsung "1inch"

### `frontend/constants/tokens.ts`
- [ ] Hapus `SOLANA_TOKENS` array seluruhnya
- [ ] Update `getTokensForChain()` — hapus `if (chain === "solana") return SOLANA_TOKENS`, langsung return `EVM_TOKENS`

### `frontend/constants/dummy-watchlist.ts`
- [ ] Hapus atau ubah dummy data yang pakai `chain: "solana"` (item id 3 dan id 4) — ganti ke `chain: "base"`

---

## 5. Contracts — Hapus Directory

| Path | Alasan |
|---|---|
| `contracts/svm/` | Seluruh directory Solana smart contracts (Anchor) — tidak dipakai lagi |

---

## 6. API Docs (Bruno) — Hapus/Edit File

| File | Action |
|---|---|
| `backend/api-docs/swap/quote-solana.bru` | Hapus — Solana swap quote test tidak relevan |
| `backend/api-docs/wallets/analyze-solana.bru` | Hapus — Solana wallet analysis test tidak relevan |
| `backend/api-docs/collection.bru` | Edit — hapus variable `solanaAddress` |

---

## 7. Urutan Eksekusi

Lakukan cleanup dalam urutan ini untuk menghindari build errors:

1. **Backend files dulu** — hapus 4 file (solana.go, birdeye.go, jupiter.go, interfaces/birdeye.go)
2. **Backend edits** — update semua file yang reference file yang dihapus (container, routes, services, config, constants)
3. **Verify backend compiles** — `cd backend && go build ./...`
4. **Frontend edits** — update constants, components, pages
5. **Verify frontend builds** — `cd frontend && npm run build`
6. **Hapus contracts/svm/** — seluruh directory
7. **Hapus API docs** — Bruno files yang tidak relevan
8. **Hapus frontend/public/chains/solana.svg**
9. **Final check** — grep seluruh codebase untuk "solana", "svm", "birdeye", "jupiter" — harus 0 results (kecuali node_modules)
