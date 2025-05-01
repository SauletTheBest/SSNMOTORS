// config/config.go
package config

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Ctx         context.Context
    Port        string
    MongoURI    string
    MongoDBName string
    JWTSecret   string
}

func Load() *Config {
    // Попытка загрузить .env из рабочей директории
    if err := godotenv.Load("../.env"); err != nil {
        log.Println("⚠️ .env file not found, using environment variables")
    }

    return &Config{
        Ctx:         context.Background(),
        Port:        getEnv("PORT", "50051"),
        MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
        MongoDBName: getEnv("MONGO_DB", "users_db"),         // NEW
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
