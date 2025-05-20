package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"inventory-service/internal/model"
	"inventory-service/internal/repository"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductUsecase struct {
	repo  repository.ProductRepository
	cache *redis.Client
}

func NewProductUsecase(repo repository.ProductRepository, cache *redis.Client) *ProductUsecase {
	return &ProductUsecase{repo: repo, cache: cache}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, p *model.Product) (string, error) {
	if p.Name == "" {
		return "", errors.New("name is required")
	}
	if p.Description == "" {
		return "", errors.New("description is required")
	}
	if p.Category == "" {
		return "", errors.New("category is required")
	}
	if p.Stock < 0 {
		return "", errors.New("stock cannot be negative")
	}
	if p.Price < 0 {
		return "", errors.New("price cannot be negative")
	}
	return u.repo.Create(ctx, p)
}

func (u *ProductUsecase) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Check Redis cache first
	cacheKey := "product:" + id
	cachedProduct, err := u.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var product model.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			log.Println("Product from cache")
			return &product, nil
		}
	}

	// If not found in cache, query MongoDB
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result in Redis (for 1 hour)
	productJSON, err := json.Marshal(product)
	if err == nil {
		u.cache.Set(ctx, cacheKey, productJSON, 5*time.Minute) // 1 hour TTL
	}
	log.Println("Prodcuct from mongo")
	return product, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, p *model.Product) error {
	if p.ID == "" {
		return errors.New("id is required")
	}
	// Все поля обязательны — нет частичного обновления
	if p.Name == "" || p.Description == "" || p.Category == "" {
		return errors.New("name, description and category are required")
	}
	if p.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return u.repo.Update(ctx, p)
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) ListProducts(ctx context.Context, category string, page, limit int32) ([]*model.Product, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit must be ≥ 1")
	}
	return u.repo.List(ctx, category, page, limit)
}

func (u *ProductUsecase) DecreaseStock(ctx context.Context, productID string, quantity int32) error {
	return u.repo.DecreaseStock(ctx, productID, quantity)
}
