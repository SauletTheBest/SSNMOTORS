package usecase

import (
	"AP2Assignment2/inventory-service/internal/domain"
	"context"
)

type auto_inc_Repo interface {
	Next(ctx context.Context, collection string) (uint64, error)
}
type product_Repo interface {
	Create(ctx context.Context, client domain.Product) error
	Update(ctx context.Context, filter domain.ProductFilter, update domain.ProductUpdateData) error
	GetWithFilter(ctx context.Context, filter domain.ProductFilter) (domain.Product, error)
	GetListWithFilter(ctx context.Context, filter domain.ProductFilter, page, limit int64) ([]domain.Product, int, error)
	Delete(ctx context.Context, filter domain.ProductFilter) error
}
