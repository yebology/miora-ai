package constants

// BaseSepoliaChainID is the chain ID for Base Sepolia testnet.
// Used by the EAS client for signing transactions.
const BaseSepoliaChainID = 84532

// BaseSepoliaEASScanURL is the base URL for the EAS explorer on Base Sepolia.
const BaseSepoliaEASScanURL = "https://base-sepolia.easscan.org"

// ChainConfig holds API-specific identifiers for each supported chain.
type ChainConfig struct {
	Name           string  // Display name
	AlchemyURL     string  // Alchemy RPC base URL (without API key)
	MoralisChainID string  // Moralis hex chain ID
	OneInchChainID string  // 1inch numeric chain ID
	DexScreenerID  string  // DexScreener chain identifier
	BlockTimeSec   float64 // Average seconds per block (for timestamp estimation)
}

// SupportedChains maps chain keys to their API configurations.
// Supports Ethereum mainnet and L2s (Arbitrum, Optimism, Base, Polygon).
var SupportedChains = map[string]ChainConfig{
	"ethereum": {
		Name:           "Ethereum",
		AlchemyURL:     "https://eth-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0x1",
		OneInchChainID: "1",
		DexScreenerID:  "ethereum",
		BlockTimeSec:   12,
	},
	"arbitrum": {
		Name:           "Arbitrum",
		AlchemyURL:     "https://arb-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0xa4b1",
		OneInchChainID: "42161",
		DexScreenerID:  "arbitrum",
		BlockTimeSec:   0.25,
	},
	"optimism": {
		Name:           "Optimism",
		AlchemyURL:     "https://opt-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0xa",
		OneInchChainID: "10",
		DexScreenerID:  "optimism",
		BlockTimeSec:   2,
	},
	"base": {
		Name:           "Base",
		AlchemyURL:     "https://base-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0x2105",
		OneInchChainID: "8453",
		DexScreenerID:  "base",
		BlockTimeSec:   2,
	},
	"polygon": {
		Name:           "Polygon",
		AlchemyURL:     "https://polygon-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0x89",
		OneInchChainID: "137",
		DexScreenerID:  "polygon",
		BlockTimeSec:   2,
	},
}

// GetChainConfig returns the config for a chain key.
// Returns nil if the chain is not supported.
func GetChainConfig(chain string) *ChainConfig {

	// Map legacy key
	if chain == "evm" {
		chain = "ethereum"
	}

	if cfg, ok := SupportedChains[chain]; ok {
		return &cfg
	}
	return nil

}

// IsEVM returns true if the chain is an EVM-compatible chain.
func IsEVM(chain string) bool {

	switch chain {
	case "evm", "ethereum", "arbitrum", "optimism", "base", "polygon":
		return true
	default:
		return false
	}

}
