package handler

import (
    "context"
    "inventory-service/internal/model"
    "inventory-service/internal/pb"
    "inventory-service/internal/usecase"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "strings"
)

type ProductHandler struct {
    pb.UnimplementedInventoryServiceServer
    uc *usecase.ProductUsecase
}

func NewProductHandler(uc *usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{uc: uc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
    prod := &model.Product{
        Name:        req.Name,
        Description: req.Description,
        Category:    req.Category,
        Stock:       req.Stock,
        Price:       req.Price,
    }
    id, err := h.uc.CreateProduct(ctx, prod)
    if err != nil {
        return nil, mapError(err)
    }
    prod.ID = id
    return &pb.ProductResponse{Product: toProto(prod)}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    prod, err := h.uc.GetProduct(ctx, req.Id)
    if err != nil {
        return nil, mapError(err)
    }
    return &pb.ProductResponse{Product: toProto(prod)}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    prod := &model.Product{
        ID:          req.Id,
        Name:        req.Name,
        Description: req.Description,
        Category:    req.Category,
        Stock:       req.Stock,
        Price:       req.Price,
    }
    if err := h.uc.UpdateProduct(ctx, prod); err != nil {
        return nil, mapError(err)
    }
    updated, _ := h.uc.GetProduct(ctx, req.Id)
    return &pb.ProductResponse{Product: toProto(updated)}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
    if err := h.uc.DeleteProduct(ctx, req.Id); err != nil {
        return nil, mapError(err)
    }
    return &pb.DeleteProductResponse{Message: "product deleted"}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    list, err := h.uc.ListProducts(ctx, req.Category, req.Page, req.Limit)
    if err != nil {
        return nil, mapError(err)
    }
    var out []*pb.Product
    for _, p := range list {
        out = append(out, toProto(p))
    }
    return &pb.ListProductsResponse{Products: out}, nil
}

func toProto(p *model.Product) *pb.Product {
    return &pb.Product{
        Id:          p.ID,
        Name:        p.Name,
        Description: p.Description,
        Category:    p.Category,
        Stock:       p.Stock,
        Price:       p.Price,
    }
}

func mapError(err error) error {
    msg := err.Error()
    switch {
    case strings.Contains(msg, "required"), strings.Contains(msg, "must be"):
        return status.Error(codes.InvalidArgument, msg)
    case strings.Contains(msg, "not found"):
        return status.Error(codes.NotFound, msg)
    case strings.Contains(msg, "already exists"):
        return status.Error(codes.AlreadyExists, msg)
    default:
        return status.Error(codes.Internal, msg)
    }
}
