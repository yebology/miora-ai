// Package config handles environment variable loading and validation.
// All configuration is read from a .env file at the project root.
// No fallback values are used — every required key must be explicitly set.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	AppPort        string // HTTP server port (APP_PORT)
	DBHost         string // PostgreSQL host (DB_HOST)
	DBPort         string // PostgreSQL port (DB_PORT)
	DBUser         string // PostgreSQL user (POSTGRES_USER)
	DBPassword     string // PostgreSQL password (POSTGRES_PASSWORD)
	DBName         string // PostgreSQL database name (POSTGRES_DB)
	AlchemyAPIKey  string // Alchemy API key for EVM RPC (ALCHEMY_API_KEY)
	MoralisAPIKey  string // Moralis API key for historical token prices (MORALIS_API_KEY)
	GeminiAPIKey   string // Google Gemini API key for AI insights (GEMINI_API_KEY)
	AllowedOrigins string // Comma-separated CORS allowed origins (ALLOWED_ORIGINS)
	Scoring        ScoringConfig
	EAS            EASConfig
	Agent          AgentConfig
}

// EASConfig holds configuration for Ethereum Attestation Service on Base Sepolia.
type EASConfig struct {
	RPCURL             string // Base Sepolia RPC endpoint (EAS_RPC_URL)
	EASContractAddress string // EAS contract address on Base Sepolia (EAS_CONTRACT_ADDRESS)
	SchemaRegistryAddr string // SchemaRegistry contract address (EAS_SCHEMA_REGISTRY_ADDRESS)
	SchemaUID          string // Registered schema UID (EAS_SCHEMA_UID)
	AttesterPrivateKey string // Private key for signing attestations (EAS_ATTESTER_PRIVATE_KEY)
}

// AgentConfig holds configuration for the AI trading agent.
type AgentConfig struct {
	CDPAPIKeyID     string // Coinbase Developer Platform API key ID (CDP_API_KEY_ID)
	CDPAPIKeySecret string // CDP API key secret (CDP_API_KEY_SECRET)
}

// ScoringConfig holds configurable thresholds for wallet scoring.
// All values are loaded from .env so they can be tuned without code changes.
type ScoringConfig struct {
	LiquidityThreshold  float64 // Tokens below this USD liquidity are considered risky (SCORING_LIQUIDITY_THRESHOLD)
	EntryTimingMaxAge   float64 // Max pair age in hours for entry timing score (SCORING_ENTRY_TIMING_MAX_AGE)
	TokenQualityLogBase float64 // Log base for token quality score, e.g. 7 means $10M = score 100 (SCORING_TOKEN_QUALITY_LOG_BASE)
}

// LoadConfig reads the .env file and validates that all required keys are present.
// Returns an error if the .env file is missing or any required key is empty.
func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	required := []string{
		"APP_PORT", "DB_HOST", "DB_PORT",
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"ALCHEMY_API_KEY", "MORALIS_API_KEY", "GEMINI_API_KEY", "ALLOWED_ORIGINS",
		"SCORING_LIQUIDITY_THRESHOLD", "SCORING_ENTRY_TIMING_MAX_AGE", "SCORING_TOKEN_QUALITY_LOG_BASE",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("missing required env: %s", key)
		}
	}

	return &Config{
		AppPort:        os.Getenv("APP_PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("POSTGRES_USER"),
		DBPassword:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		AlchemyAPIKey:  os.Getenv("ALCHEMY_API_KEY"),
		MoralisAPIKey:  os.Getenv("MORALIS_API_KEY"),
		GeminiAPIKey:   os.Getenv("GEMINI_API_KEY"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		Scoring: ScoringConfig{
			LiquidityThreshold:  getEnvFloat("SCORING_LIQUIDITY_THRESHOLD"),
			EntryTimingMaxAge:   getEnvFloat("SCORING_ENTRY_TIMING_MAX_AGE"),
			TokenQualityLogBase: getEnvFloat("SCORING_TOKEN_QUALITY_LOG_BASE"),
		},
		EAS: EASConfig{
			RPCURL:             os.Getenv("EAS_RPC_URL"),
			EASContractAddress: getEnvDefault("EAS_CONTRACT_ADDRESS", "0x4200000000000000000000000000000000000021"),
			SchemaRegistryAddr: getEnvDefault("EAS_SCHEMA_REGISTRY_ADDRESS", "0x4200000000000000000000000000000000000020"),
			SchemaUID:          os.Getenv("EAS_SCHEMA_UID"),
			AttesterPrivateKey: os.Getenv("EAS_ATTESTER_PRIVATE_KEY"),
		},
		Agent: AgentConfig{
			CDPAPIKeyID:     os.Getenv("CDP_API_KEY_ID"),
			CDPAPIKeySecret: os.Getenv("CDP_API_KEY_SECRET"),
		},
	}, nil

}

// DSN returns the PostgreSQL connection string in key=value format.
// Used by GORM to establish the database connection.
func (c *Config) DSN() string {

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)

}

// getEnvFloat parses an environment variable as float64. Returns 0 if invalid.
func getEnvFloat(key string) float64 {

	val, _ := strconv.ParseFloat(os.Getenv(key), 64)
	return val

}

// getEnvDefault returns the env value or a default string if not set.
func getEnvDefault(key, defaultVal string) string {

	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val

}

// getEnvFloatDefault parses an environment variable as float64 with a fallback default.
func getEnvFloatDefault(key string, defaultVal float64) float64 {

	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return defaultVal
	}
	return val

}
