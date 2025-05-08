package handler

import (
	"AP2Assignment2/inventory-service/internal/domain"
	"context"
)

type ProductUseCase interface {
	Create(ctx context.Context, product domain.Product) (domain.Product, error)

	//Get gets product with specified filter
	Get(ctx context.Context, pf domain.ProductFilter) (domain.Product, error)

	//GetAll returns list of products(in the specified page with the specified limit) and total amount of products
	GetAll(ctx context.Context, pf domain.ProductFilter, page, limit int64) ([]domain.Product, int, error)

	Update(ctx context.Context, filter domain.ProductFilter, updated domain.ProductUpdateData) error

	Delete(ctx context.Context, filter domain.ProductFilter) error
}
