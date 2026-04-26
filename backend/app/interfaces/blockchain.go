package interfaces

type TransferData struct {
	Hash            string
	From            string
	To              string
	Value           string
	TokenSymbol     string
	ContractAddress string // Token contract address (from rawContract.address in Alchemy)
	BlockNumber     uint64
	Timestamp       int64
	Direction       string // "in" = received, "out" = sent
}

type BlockchainClient interface {
	GetTransfers(address string, limit int, chain ...string) ([]TransferData, error)
}
