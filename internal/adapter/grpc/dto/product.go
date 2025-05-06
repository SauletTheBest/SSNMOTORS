package dto

import (
	"AP2Assignment2/inventory-service/internal/domain"
	proto "AP2Assignment2/inventory-service/protos/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CreateProductRequest struct {
	Name     string
	Category string
	Price    float64
	Stock    uint64
}

type ProductResponse struct {
	ID        uint64
	Name      string
	Category  string
	Price     float64
	Stock     uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetProductRequest struct {
	ProductID uint64
}

type UpdateProductRequest struct {
	ProductID uint64
	Name      *string
	Category  *string
	Price     *float64
	Stock     *uint64
}

type ListProductsRequest struct {
	Name     *string
	Category *string
	Price    *float64
	Stock    *uint64
	Page     int64
	Limit    int64
}

type DeleteProductRequest struct {
	ProductID uint64
}

// FromCreateRequestProto converts gRPC request to DTO
func FromCreateRequestProto(req *proto.CreateProductRequest) *CreateProductRequest {
	return &CreateProductRequest{
		Name:     req.Name,
		Category: req.Category,
		Price:    req.Price,
		Stock:    req.Stock,
	}
}

// ToProduct converts DTO to domain model
func (d *CreateProductRequest) ToProduct() domain.Product {
	return domain.Product{
		Name:     d.Name,
		Category: d.Category,
		Price:    d.Price,
		Stock:    d.Stock,
	}
}

// FromProduct converts domain model to DTO
func FromProduct(product domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Category:  product.Category,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

// ToProtoProductResponse converts DTO to gRPC response
func (d *ProductResponse) ToProtoProductResponse() *proto.ProductResponse {
	return &proto.ProductResponse{
		ProductId: d.ID,
		Name:      d.Name,
		Category:  d.Category,
		Price:     d.Price,
		Stock:     d.Stock,
		CreatedAt: timestamppb.New(d.CreatedAt).String(),
		UpdatedAt: timestamppb.New(d.UpdatedAt).String(),
	}
}

// FromGetRequestProto converts gRPC request to DTO
func FromGetRequestProto(req *proto.GetProductRequest) *GetProductRequest {
	return &GetProductRequest{
		ProductID: req.ProductId,
	}
}

// ToDomainFilter converts DTO to domain filter
func (d *GetProductRequest) ToDomainFilter() domain.ProductFilter {
	return domain.ProductFilter{
		ID: &d.ProductID,
	}
}

// FromUpdateRequestProto converts gRPC request to DTO
func FromUpdateRequestProto(req *proto.UpdateProductRequest) *UpdateProductRequest {
	return &UpdateProductRequest{
		ProductID: req.ProductId,
		Name:      req.Name,
		Category:  req.Category,
		Price:     req.Price,
		Stock:     req.Stock,
	}
}

// ToDomainFilterAndUpdate converts DTO to domain filter and update data
func (d *UpdateProductRequest) ToDomainFilterAndUpdate() (domain.ProductFilter, domain.ProductUpdateData) {
	filter := domain.ProductFilter{
		ID: &d.ProductID,
	}
	update := domain.ProductUpdateData{
		Name:     d.Name,
		Category: d.Category,
		Price:    d.Price,
		Stock:    d.Stock,
	}
	return filter, update
}

// FromListRequestProto converts gRPC request to DTO
func FromListRequestProto(req *proto.ListProductsRequest) *ListProductsRequest {
	return &ListProductsRequest{
		Name:     req.Name,
		Category: req.Category,
		Price:    req.Price,
		Stock:    req.Stock,
		Page:     req.Page,
		Limit:    req.Limit,
	}
}

// ToDomainFilter converts DTO to domain filter
func (d *ListProductsRequest) ToDomainFilter() domain.ProductFilter {
	return domain.ProductFilter{
		Name:     d.Name,
		Category: d.Category,
		Price:    d.Price,
		Stock:    d.Stock,
	}
}

// FromDeleteRequestProto converts gRPC request to DTO
func FromDeleteRequestProto(req *proto.DeleteProductRequest) *DeleteProductRequest {
	return &DeleteProductRequest{
		ProductID: req.ProductId,
	}
}

// ToDomainFilter converts DTO to domain filter
func (d *DeleteProductRequest) ToDomainFilter() domain.ProductFilter {
	return domain.ProductFilter{
		ID: &d.ProductID,
	}
}
