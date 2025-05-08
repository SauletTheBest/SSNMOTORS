package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"api-gateway/config"
	"api-gateway/internal/pb/inventory"
	"api-gateway/internal/pb/order"
	"api-gateway/internal/pb/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Handler manages REST handlers and gRPC clients
type Handler struct {
	inventoryClient inventory.InventoryServiceClient
	orderClient     order.OrderServiceClient
	userClient      user.UserServiceClient
}

// NewHandler initializes gRPC clients and returns a Handler
func NewHandler(cfg *config.Config) (*Handler, error) {
	// Connect to inventory service
	inventoryConn, err := grpc.Dial(cfg.InventoryService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to order service
	orderConn, err := grpc.Dial(cfg.OrderService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to user service
	userConn, err := grpc.Dial(cfg.UserService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Handler{
		inventoryClient: inventory.NewInventoryServiceClient(inventoryConn),
		orderClient:     order.NewOrderServiceClient(orderConn),
		userClient:      user.NewUserServiceClient(userConn),
	}, nil
}

// Inventory Handlers
func (h *Handler) CreateProduct(c *gin.Context) {
	var req inventory.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.inventoryClient.CreateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Product)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	req := &inventory.GetProductRequest{Id: id}

	resp, err := h.inventoryClient.GetProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Product)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var req inventory.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.inventoryClient.UpdateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Product)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	req := &inventory.DeleteProductRequest{Id: id}

	resp, err := h.inventoryClient.DeleteProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (h *Handler) ListProducts(c *gin.Context) {
	category := c.Query("category")
	var page, limit int32
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = int32(val)
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = int32(val)
		}
	}

	req := &inventory.ListProductsRequest{
		Category: category,
		Page:     page,
		Limit:    limit,
	}

	resp, err := h.inventoryClient.ListProducts(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Products)
}

// Order Handlers
func (h *Handler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.orderClient.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": resp.Id, "message": resp.Message})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	req := &order.GetOrderRequest{Id: id}

	resp, err := h.orderClient.GetOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var req order.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.orderClient.UpdateOrderStatus(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": resp.Id, "message": resp.Message})
}

func (h *Handler) ListUserOrders(c *gin.Context) {
	userID := c.Query("user_id")
	req := &order.ListUserOrdersRequest{UserId: userID}

	resp, err := h.orderClient.ListUserOrders(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Orders)
}

// User Handlers
func (h *Handler) RegisterUser(c *gin.Context) {
	var req user.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userClient.RegisterUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": resp.Id, "message": resp.Message})
}

func (h *Handler) AuthenticateUser(c *gin.Context) {
	var req user.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userClient.AuthenticateUser(context.Background(), &req)
	if err != nil {
		log.Printf("AuthenticateUser failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User authenticated: username=%s, token=%s", req.Username, resp.Token)
	c.JSON(http.StatusOK, gin.H{"token": resp.Token, "message": resp.Message})
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	req := &user.UserID{Id: id}

	resp, err := h.userClient.GetUserProfile(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}