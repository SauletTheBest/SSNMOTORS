package config

import "os"

// Config содержит настройки подключения к микросервисам и HTTP-порт
type Config struct {
    UserServiceAddr     string // Адрес user-service
    InventoryServiceAddr string // Адрес inventory-service
    OrderServiceAddr    string // Адрес order-service
    HttpPort            string // Порт для HTTP-сервера
}

// Load загружает конфигурацию из переменных окружения или использует значения по умолчанию
func Load() *Config {
    return &Config{
        UserServiceAddr:     getEnv("USER_SERVICE_ADDR", "localhost:5051"),
        InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:5052"),
        OrderServiceAddr:    getEnv("ORDER_SERVICE_ADDR", "localhost:5053"),
        HttpPort:            getEnv("HTTP_PORT", ":8080"),
    }
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}