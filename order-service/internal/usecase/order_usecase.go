package usecase

import (
	"context"
	"errors"
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	if order.UserID == "" || len(order.Products) == 0 {
		return "", errors.New("invalid order data")
	}
	if order.Status == "" {
		order.Status = "PENDING"
	}
	return u.repo.Create(ctx, order)
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
	return u.repo.FindByUserID(ctx, userID)
}