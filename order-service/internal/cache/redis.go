// internal/cache/redis.go
package cache

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Адрес Redis
	})

	// Пингуем Redis, чтобы убедиться, что подключение установлено
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	return client
}
