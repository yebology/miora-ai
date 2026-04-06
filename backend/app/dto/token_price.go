package dto

// TokenPriceData holds token price info from Moralis.
type TokenPriceData struct {
	TokenAddress string
	UsdPrice     float64
	BlockNumber  uint64
}
