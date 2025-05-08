package handler

import (
	"AP2Assignment2/inventory-service/internal/adapter/http/handler/dto"
	"AP2Assignment2/inventory-service/internal/domain"
	"errors"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	useCase ProductUseCase
}

type IProductHandler interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

func NewProductHandler(uc ProductUseCase) *ProductHandler {
	return &ProductHandler{uc}
}

// Create handles creating a new product
func (h *ProductHandler) Create(ctx *gin.Context) {
	product, err := dto.FromProductRequest(ctx)
	if err != nil {
		return // Error response already sent in FromProductRequest to ctx
	}

	createdProduct, err := h.useCase.Create(ctx.Request.Context(), product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToProductResponse(createdProduct))
}

// GetByID handles retrieving a product by its ID
func (h *ProductHandler) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	filter := domain.ProductFilter{ID: &id}
	product, err := h.useCase.Get(ctx.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve product: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductResponse(product))
}

// GetAll handles listing products with pagination
func (h *ProductHandler) GetAll(ctx *gin.Context) {
	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	// Parse filter parameters
	var filter domain.ProductFilter
	if category := ctx.Query("category"); category != "" {
		filter.Category = &category
	}
	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}

	// Fetch products from use case
	products, total, err := h.useCase.GetAll(ctx.Request.Context(), filter, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve products: " + err.Error()})
		return
	}

	// Convert to response DTOs
	responses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = dto.ToProductResponse(product)
	}

	// Calculate total pages
	totalPages := total / int(limit)
	if total%int(limit) > 0 {
		totalPages++
	}

	// Create paginated response
	response := dto.ProductListResponse{
		Items:      responses,
		Total:      total,
		Page:       int(page),
		Limit:      int(limit),
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}

// Update handles updating a product
func (h *ProductHandler) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	// Parse update data from request body
	var updateReq struct {
		Name     *string  `json:"name"`
		Category *string  `json:"category"`
		Stock    *uint64  `json:"stock"`
		Price    *float64 `json:"price"`
	}

	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	// Create filter and update data
	filter := domain.ProductFilter{ID: &id}
	updateData := domain.ProductUpdateData{
		Name:     updateReq.Name,
		Category: updateReq.Category,
		Stock:    updateReq.Stock,
		Price:    updateReq.Price,
	}

	// Check if at least one field is being updated
	if updateReq.Name == nil && updateReq.Category == nil && updateReq.Stock == nil && updateReq.Price == nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "No fields to update"})
		return
	}

	// Update the product
	err = h.useCase.Update(ctx.Request.Context(), filter, updateData)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update product: " + err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Delete handles deleting a product
func (h *ProductHandler) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	filter := domain.ProductFilter{ID: &id}
	err = h.useCase.Delete(ctx.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete product: " + err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
