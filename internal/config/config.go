package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey       string
	HANAHost     string
	HANAPort     int
	HANAUser     string
	HANAPassword string
	HANASchema   string
}

func Load() (*Config, error) {
	// Load .env if it exists.
	// In production, environment variables can be provided directly.
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("HANA_PORT", "30015"))
	if err != nil {
		return nil, fmt.Errorf("invalid HANA_PORT: %w", err)
	}

	cfg := &Config{
		APIKey:       os.Getenv("API_KEY"),
		HANAHost:     os.Getenv("HANA_HOST"),
		HANAPort:     port,
		HANAUser:     os.Getenv("HANA_USER"),
		HANAPassword: os.Getenv("HANA_PASSWORD"),
		HANASchema:   os.Getenv("HANA_SCHEMA"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY is required")
	}

	if c.HANAHost == "" {
		return fmt.Errorf("HANA_HOST is required")
	}

	if c.HANAUser == "" {
		return fmt.Errorf("HANA_USER is required")
	}

	if c.HANASchema == "" {
		return fmt.Errorf("HANA_SCHEMA is required")
	}

	return nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
