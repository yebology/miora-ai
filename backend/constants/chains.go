// Package constants defines supported chain configurations.
//
// Miora V2 is Base-only. Other EVM chains have been removed for the hackathon.
// Multi-chain support can be re-added post-hackathon.
package constants

// BaseSepoliaChainID is the chain ID for Base Sepolia testnet.
// Used by the EAS client for signing transactions.
const BaseSepoliaChainID = 84532

// BaseSepoliaEASScanURL is the base URL for the EAS explorer on Base Sepolia.
const BaseSepoliaEASScanURL = "https://base-sepolia.easscan.org"

// ChainConfig holds API-specific identifiers for Base.
type ChainConfig struct {
	Name           string  // Display name
	AlchemyURL     string  // Alchemy RPC base URL (without API key)
	MoralisChainID string  // Moralis hex chain ID
	DexScreenerID  string  // DexScreener chain identifier
	BlockTimeSec   float64 // Average seconds per block (for timestamp estimation)
}

// SupportedChains maps chain keys to their API configurations.
// Base-only for V2 hackathon.
var SupportedChains = map[string]ChainConfig{
	"base": {
		Name:           "Base",
		AlchemyURL:     "https://base-mainnet.g.alchemy.com/v2/",
		MoralisChainID: "0x2105",
		DexScreenerID:  "base",
		BlockTimeSec:   2,
	},
}

// GetChainConfig returns the config for a chain key.
// Returns nil if the chain is not supported.
func GetChainConfig(chain string) *ChainConfig {
	if chain == "evm" {
		chain = "base"
	}

	if cfg, ok := SupportedChains[chain]; ok {
		return &cfg
	}
	return nil
}

// IsEVM returns true if the chain is supported (Base only for V2).
func IsEVM(chain string) bool {
	switch chain {
	case "evm", "base":
		return true
	default:
		return false
	}
}
