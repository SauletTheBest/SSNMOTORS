package repository

import (
    "context"
    "inventory-service/internal/model"
)

type ProductRepository interface {
    Create(ctx context.Context, p *model.Product) (string, error)
    GetByID(ctx context.Context, id string) (*model.Product, error)
    Update(ctx context.Context, p *model.Product) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, category string, page, limit int32) ([]*model.Product, error)
    DecreaseStock(ctx context.Context, productID string, quantity int32) error
}
