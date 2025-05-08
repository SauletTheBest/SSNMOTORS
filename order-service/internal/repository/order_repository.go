package repository

import (
	"context"
	"order-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	FindByID(ctx context.Context, id string) (*model.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
