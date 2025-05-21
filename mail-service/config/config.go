package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MailerSendAPIKey string
	ServerPort       string
}

func LoadConfig() *Config {
	// Загружаем .env файл
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️ .env файл не найден, используем переменные окружения")
	}

	config := &Config{
		MailerSendAPIKey: getEnv("MAILERSEND_API_KEY", ""),
		ServerPort:       getEnv("SERVER_PORT", "50054"),
	}

	if config.MailerSendAPIKey == "" {
		log.Fatal("❌ MAILERSEND_API_KEY не задан")
	}

	return config
}

// Вспомогательная функция с дефолтным значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
