package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"encoding/json"
	"order-service/internal/model"
	"order-service/internal/queue"
	"order-service/internal/repository"

	"github.com/redis/go-redis/v9"
)

type OrderUsecase struct {
	repo      repository.OrderRepository
	publisher queue.Publisher
	cache     *redis.Client
}

func NewOrderUsecase(repo repository.OrderRepository, publisher queue.Publisher, cache *redis.Client) *OrderUsecase {
	return &OrderUsecase{
		repo:      repo,
		publisher: publisher,
		cache:     cache,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	if order.UserID == "" || len(order.Products) == 0 {
		return "", errors.New("invalid order data")
	}
	if order.Status == "" {
		order.Status = "PENDING"
	}

	// Basic validation for ProductID (ensure it's not a name, but this needs more robust validation)
	for _, product := range order.Products {
		if product.ProductID == "" {
			return "", errors.New("invalid product ID in order")
		}
		// Ideally, validate that ProductID looks like a MongoDB ObjectID (e.g., 24-character hex string)
		if len(product.ProductID) != 24 {
			log.Printf("Warning: ProductID %s does not match expected ObjectID format", product.ProductID)
		}
	}

	// Save the order to the database
	id, err := u.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	order.ID = id

	// Publish to NATS
	err = u.publisher.PublishOrderCreated(ctx, order)
	if err != nil {
		log.Printf("[NATS] Failed to publish order.created event: %v", err)

		// Log the error but don't fail the request (fire-and-forget)
	}

	return id, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	if id == "" {
		return nil, errors.New("invalid order ID")
	}
	return u.repo.FindByID(ctx, id)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	if id == "" || status == "" {
		return errors.New("invalid input data")
	}
	validStatuses := map[string]bool{
		"PENDING":   true,
		"COMPLETED": true,
		"CANCELLED": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *OrderUsecase) ListUserOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}

	// Проверяем кэш перед запросом в базу данных
	cacheKey := "user_orders:" + userID
	cachedOrders, err := u.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var orders []*model.Order
		if err := json.Unmarshal([]byte(cachedOrders), &orders); err == nil {
			log.Println("Orders from cache")
			return orders, nil
		}
	}

	// Если данных нет в кэше, запрашиваем из базы данных
	orders, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем результат в кэш
	ordersJson, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	// Указываем TTL в 5 минут (300 секунд)

	err = u.cache.Set(ctx, cacheKey, ordersJson, 5*time.Minute).Err()
	if err != nil {
		log.Println("Failed to cache orders:", err)
	}

	log.Println("Orders from mongo")
	return orders, nil
}
