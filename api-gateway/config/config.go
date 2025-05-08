package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port             string
	InventoryService string
	OrderService     string
	UserService      string
	JWTSecret        string
}

// Load loads configuration from environment variables or .env file
func Load() (*Config, error) {
	// Try to load .env file, continue if not found
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	} else {
		log.Println("Using configuration from .env file")
	}

	// Create config with values from environment
	cfg := &Config{
		Port:             getEnvWithDefault("PORT", "8080"),
		InventoryService: getEnvWithDefault("INVENTORY_SERVICE", "localhost:50051"),
		OrderService:     getEnvWithDefault("ORDER_SERVICE", "localhost:50052"),
		UserService:      getEnvWithDefault("USER_SERVICE", "localhost:50053"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
	}

	// Validate JWT secret
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

// getEnvWithDefault retrieves an environment variable or returns a default value if not set
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}