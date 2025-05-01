package main

import (
    "api-gateway/config"
    "api-gateway/internal/server"
)

func main() {
    // Загрузка конфигурации
    cfg := config.Load()

    // Создание и запуск сервера
    server := server.NewServer(cfg)
    server.Start()
}