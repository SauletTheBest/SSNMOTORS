package repository

import (
	"context"
	"order-service/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	FindByID(ctx context.Context, id string) (*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	FindByUserID(ctx context.Context, userID string) ([]*model.Order, error)
}