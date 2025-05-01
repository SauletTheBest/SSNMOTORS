package handler

import (
    "api-gateway/internal/pb/inventory"
    "api-gateway/internal/pb/order"
    "api-gateway/internal/pb/user"
    "net/http" // Import the pb package
    "api-gateway/internal/middleware"
    "github.com/gin-gonic/gin"
)

// Handler содержит gRPC клиенты для взаимодействия с микросервисами
type Handler struct {
    inventoryClient inventory.InventoryServiceClient
    orderClient     order.OrderServiceClient
    userClient      user.UserServiceClient
}

func NewHandler(
    inventoryClient inventory.InventoryServiceClient,
    orderClient order.OrderServiceClient,
    userClient user.UserServiceClient,
) *Handler {
    return &Handler{
        inventoryClient: inventoryClient,
        orderClient:     orderClient,
        userClient:      userClient,
    }
}
// --- User Handlers ---

func (h *Handler) RegisterUser(c *gin.Context) {
    var req pb.UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    resp, err := h.userClient.RegisterUser(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) Login(c *gin.Context) {
    var req pb.AuthRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    resp, err := h.userClient.AuthenticateUser(c, &req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Генерация JWT токена
    token, err := middleware.GenerateToken(resp.UserId) // Replace UserId with the correct field name, e.g., resp.UserID
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
        return
    }

    // Отправляем токен клиенту
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetProfile(c *gin.Context) {
    id := c.Param("id")
    resp, err := h.userClient.GetUserProfile(c, &pb.UserID{Id: id})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

// --- Inventory Handlers ---

func (h *Handler) CreateProduct(c *gin.Context) {
    var req pb.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    resp, err := h.inventoryClient.CreateProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProduct(c *gin.Context) {
    id := c.Param("id")

    // Вызов метода микросервиса
    resp, err := h.inventoryClient.GetProduct(c, &inventory.GetProductRequest{Id: id})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // Отправка ответа клиенту
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
    var req pb.UpdateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    req.Id = id
    resp, err := h.inventoryClient.UpdateProduct(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    resp, err := h.inventoryClient.DeleteProduct(c, &pb.DeleteProductRequest{Id: id})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListProducts(c *gin.Context) {
    category := c.Query("category")
    page := int32(1) // default
    limit := int32(10)
    resp, err := h.inventoryClient.ListProducts(c, &pb.ListProductsRequest{
        Category: category,
        Page:     page,
        Limit:    limit,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

// --- Order Handlers ---

func (h *Handler) CreateOrder(c *gin.Context) {
    var req order.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Вызов метода микросервиса
    resp, err := h.orderClient.CreateOrder(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Отправка ответа клиенту
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetOrder(c *gin.Context) {
    id := c.Param("id")
    resp, err := h.orderClient.GetOrder(c, &pb.GetOrderRequest{Id: id})
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListUserOrders(c *gin.Context) {
    id := c.Param("id")
    resp, err := h.orderClient.ListUserOrders(c, &pb.ListUserOrdersRequest{UserId: id})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
    id := c.Param("id")
    var req pb.UpdateOrderStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    req.Id = id
    resp, err := h.orderClient.UpdateOrderStatus(c, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, resp)
}