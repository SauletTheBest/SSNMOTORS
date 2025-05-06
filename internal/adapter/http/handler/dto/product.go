package dto

import (
	"AP2Assignment2/inventory-service/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// ProductRequest represents the request body for creating a product
type ProductRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Stock    uint64  `json:"stock" binding:"required,min=0"`
	Price    float64 `json:"price" binding:"required,min=0"`
}

// ProductResponse represents the response body after creating a product
type ProductResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Stock     uint64    `json:"stock"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// FromProductRequest converts a Gin request to a Product model
func FromProductRequest(ctx *gin.Context) (domain.Product, error) {
	var req ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return domain.Product{}, err
	}

	return domain.Product{
		Name:     req.Name,
		Category: req.Category,
		Stock:    req.Stock,
		Price:    req.Price,
	}, nil
}

// ToProductResponse converts a Product model to a response DTO
func ToProductResponse(product domain.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Category:  product.Category,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
