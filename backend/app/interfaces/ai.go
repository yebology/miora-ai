package interfaces

// IAI defines the contract for AI text generation.
type IAI interface {
	Generate(prompt string) (string, error)
}
