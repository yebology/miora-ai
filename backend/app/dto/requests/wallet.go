package requests

type AnalyzeWallet struct {
	Address string `json:"address" validate:"required"`
	Chain   string `json:"chain" validate:"required,oneof=evm svm"`
}
