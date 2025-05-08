package config

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Ctx         context.Context
    MongoURI    string
    MongoDBName string
    Port        string
    NATSURL     string
}

func Load() *Config {
    // Попытка загрузить .env файл
    if err := godotenv.Load("../.env"); err != nil {
        log.Println("No .env file found, using environment variables directly")
    }

    // Получение переменных окружения
    return &Config{
        Ctx:         context.TODO(),
        MongoURI:    getEnv("MONGO_URI"),
        MongoDBName: getEnv("MONGO_DB"),
        Port:        getEnv("PORT"),
        NATSURL:    getEnv("NATS_URL"),
    }
}

func getEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("Environment variable %s is not set", key)
    }
    return value
}