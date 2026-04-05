package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	AlchemyAPIKey string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	required := []string{
		"APP_PORT", "DB_HOST", "DB_PORT",
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"ALCHEMY_API_KEY",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("missing required env: %s", key)
		}
	}

	return &Config{
		AppPort:       os.Getenv("APP_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("POSTGRES_USER"),
		DBPassword:    os.Getenv("POSTGRES_PASSWORD"),
		DBName:        os.Getenv("POSTGRES_DB"),
		AlchemyAPIKey: os.Getenv("ALCHEMY_API_KEY"),
	}, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
