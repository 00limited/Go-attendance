package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv() error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
			return err
		}
	}
	return nil
}

// GetEnv retrieves an environment variable or returns a default value if not set
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
