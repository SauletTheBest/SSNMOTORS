package config

import (
	"log"
	"os"
    "github.com/joho/godotenv"
)

type Config struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	cfg := &Config{
		SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASSWORD"),
	}

	if cfg.SMTPUser == "" || cfg.SMTPPass == "" {
		log.Fatal("SMTP_USER and SMTP_PASSWORD must be set in environment variables")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
