package interfaces

type TransferData struct {
	Hash        string
	From        string
	To          string
	Value       string
	TokenSymbol string
	BlockNumber uint64
	Timestamp   int64
}

type BlockchainClient interface {
	GetTransfers(address string) ([]TransferData, error)
}
