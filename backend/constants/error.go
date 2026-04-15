package constants

// Error messages used in API responses via output.GetError().
const (
	DataNotFound     = "We couldn't find the data you're looking for."
	InvalidData      = "The data you provided is invalid."
	InvalidRequest   = "Invalid request body."
	UnsupportedChain = "Unsupported chain."
	AddressRequired  = "Address and chain are required."
	AnalysisFailed   = "Failed to analyze wallet."
	Unauthorized     = "Unauthorized. Please provide a valid token."
)

// HTTP status codes not defined in fiber.
const (
	HTTPStatusPaymentRequired = 402 // x402 payment required
)

// x402 protocol constants.
const (
	X402FacilitatorURL = "https://facilitator.x402.rs" // Public x402 payment facilitator
)
