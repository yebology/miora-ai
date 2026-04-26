package dto

// TokenPairData holds token pair info from DexScreener, used across layers.
type TokenPairData struct {
	ChainID        string
	PairAddress    string
	PriceUsd       string
	Liquidity      float64
	VolumeH24      float64
	FDV            float64
	MarketCap      float64
	PairCreatedAt  int64
	PriceChangeH24 float64
	BaseSymbol     string
	BaseAddress    string
}
