package handler

import (
	"context"
	"order-service/internal/pb"
	"order-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

// CreateOrder handles creating a new order.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	if req.UserId == "" || req.CarId == "" || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing or invalid required fields")
	}

	order, err := h.uc.Create(ctx, req.UserId, req.CarId, req.Quantity)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create order: "+err.Error())
	}

	return &pb.OrderResponse{
		Id:        order.ID,
		UserId:    order.UserID,
		CarId:     order.CarID,
		Quantity:  order.Quantity,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetOrder returns an order by its ID.
func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.OrderIdRequest) (*pb.OrderResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "Order ID is required")
	}

	order, err := h.uc.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "order not found: "+err.Error())
	}

	return &pb.OrderResponse{
		Id:        order.ID,
		UserId:    order.UserID,
		CarId:     order.CarID,
		Quantity:  order.Quantity,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetOrdersByUser returns all orders made by a specific user.
func (h *OrderHandler) GetOrdersByUser(ctx context.Context, req *pb.UserIdRequest) (*pb.OrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	orders, err := h.uc.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get orders: "+err.Error())
	}

	var pbOrders []*pb.OrderResponse
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.OrderResponse{
			Id:        order.ID,
			UserId:    order.UserID,
			CarId:     order.CarID,
			Quantity:  order.Quantity,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.OrdersResponse{Orders: pbOrders}, nil
}

// UpdateOrderStatus changes the status of an existing order.
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	if req.OrderId == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "Order ID and new status are required")
	}

	err := h.uc.UpdateStatus(ctx, req.OrderId, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update order status: "+err.Error())
	}

	return &pb.UpdateOrderResponse{
		Id:      req.OrderId,
		Message: "Order status updated successfully",
	}, nil
}

// CancelOrder sets the status of an order to "cancelled".
func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.OrderIdRequest) (*pb.OrderResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "Order ID is required")
	}

	order, err := h.uc.Cancel(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to cancel order: "+err.Error())
	}

	return &pb.OrderResponse{
		Id:        order.ID,
		UserId:    order.UserID,
		CarId:     order.CarID,
		Quantity:  order.Quantity,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
