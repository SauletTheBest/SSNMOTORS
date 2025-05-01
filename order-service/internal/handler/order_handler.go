package handler

import (
	"context"
	"order-service/internal/model"
	"order-service/internal/pb"
	"order-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(u *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &model.Order{
		UserID: req.UserId,
		Total:  req.Total,
	}
	for _, p := range req.Products {
		order.Products = append(order.Products, model.Product{
			ProductID: p.ProductId,
			Quantity:  int(p.Quantity),
		})
	}
	id, err := h.usecase.CreateOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.CreateOrderResponse{
		Id:      id,
		Message: "Order created successfully",
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.usecase.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	resp := &pb.GetOrderResponse{
		Id:     order.ID,
		UserId: order.UserID,
		Total:  order.Total,
		Status: order.Status,
	}
	for _, p := range order.Products {
		resp.Products = append(resp.Products, &pb.Product{
			ProductId: p.ProductID,
			Quantity:  int32(p.Quantity),
		})
	}
	return resp, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	err := h.usecase.UpdateOrderStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.UpdateOrderStatusResponse{
		Id:      req.Id,
		Message: "Order status updated successfully",
	}, nil
}

func (h *OrderHandler) ListUserOrders(ctx context.Context, req *pb.ListUserOrdersRequest) (*pb.ListUserOrdersResponse, error) {
	orders, err := h.usecase.ListUserOrders(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	resp := &pb.ListUserOrdersResponse{}
	for _, order := range orders {
		orderResp := &pb.GetOrderResponse{
			Id:     order.ID,
			UserId: order.UserID,
			Total:  order.Total,
			Status: order.Status,
		}
		for _, p := range order.Products {
			orderResp.Products = append(orderResp.Products, &pb.Product{
				ProductId: p.ProductID,
				Quantity:  int32(p.Quantity),
			})
		}
		resp.Orders = append(resp.Orders, orderResp)
	}
	return resp, nil
}