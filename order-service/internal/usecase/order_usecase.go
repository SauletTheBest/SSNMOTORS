package usecase

import (
	"context"
	"errors"
	"order-service/internal/model"
	"order-service/internal/repository"
	"time"
)

type OrderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

// CreateOrder creates a new order
func (u *OrderUsecase) Create(ctx context.Context, userID, carID string, quantity int32) (*model.Order, error) {
	if userID == "" || carID == "" || quantity <= 0 {
		return nil, errors.New("missing or invalid fields")
	}

	order := &model.Order{
		UserID:    userID,
		CarID:     carID,
		Quantity:  quantity,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	orderID, err := u.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	order.ID = orderID
	return order, nil
}

// GetOrderByID returns order by ID
func (u *OrderUsecase) GetByID(ctx context.Context, id string) (*model.Order, error) {
	if id == "" {
		return nil, errors.New("order ID required")
	}
	return u.repo.FindByID(ctx, id)
}

// ListOrdersByUser returns all orders by a user
func (u *OrderUsecase) GetByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID required")
	}
	return u.repo.FindByUserID(ctx, userID)
}

// UpdateStatus updates the status of an order
func (u *OrderUsecase) UpdateStatus(ctx context.Context, orderID, status string) error {
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == status {
		return errors.New("order already in the requested status")
	}

	order.Status = status
	return u.repo.Update(ctx, order)
}

// CancelOrder sets status to "cancelled"
func (u *OrderUsecase) Cancel(ctx context.Context, id string) (*model.Order, error) {
	order, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order.Status == "cancelled" {
		return nil, errors.New("order already cancelled")
	}

	order.Status = "cancelled"

	if err := u.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
