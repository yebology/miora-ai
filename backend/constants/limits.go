package constants

// TransactionLimitConfig holds the available fetch limits and default selection
// for transaction history queries per chain type.
type TransactionLimitConfig struct {
	Options []int // Available limit options for the user to choose from
	Default int   // Default limit (first option shown / pre-selected)
}

// EVMTransactionLimits defines fetch limits for EVM chains (Ethereum + L2s).
var EVMTransactionLimits = TransactionLimitConfig{
	Options: []int{10, 25, 50, 100},
	Default: 10,
}

// SVMTransactionLimits defines fetch limits for Solana.
var SVMTransactionLimits = TransactionLimitConfig{
	Options: []int{20, 50, 100, 200},
	Default: 20,
}

// GetTransactionLimits returns the transaction limit config based on chain key.
func GetTransactionLimits(chain string) TransactionLimitConfig {
	if IsSolana(chain) {
		return SVMTransactionLimits
	}
	return EVMTransactionLimits
}

// IsValidTransactionLimit checks if the given limit is a valid option for the chain.
func IsValidTransactionLimit(chain string, limit int) bool {
	cfg := GetTransactionLimits(chain)
	for _, opt := range cfg.Options {
		if opt == limit {
			return true
		}
	}
	return false
}
