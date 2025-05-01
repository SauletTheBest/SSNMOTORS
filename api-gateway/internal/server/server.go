package server

import (
    "api-gateway/config"
    "api-gateway/internal/handler"
    "api-gateway/internal/pb/inventory"
	"api-gateway/internal/pb/order"
	"api-gateway/internal/pb/user"
    "context"	
    "log"
	"api-gateway/internal/middleware"
    "time"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
)

type Server struct {
    cfg            *config.Config
    inventoryClient inventory.InventoryServiceClient
    orderClient    order.OrderServiceClient
    userClient     user.UserServiceClient
}

func NewServer(cfg *config.Config) *Server {
    inventoryConn := mustConnect(cfg.InventoryServiceAddr)
    orderConn := mustConnect(cfg.OrderServiceAddr)
    userConn := mustConnect(cfg.UserServiceAddr)

    return &Server{
        cfg: cfg,
        inventoryClient: inventory.NewInventoryServiceClient(inventoryConn),
        orderClient:     order.NewOrderServiceClient(orderConn),
        userClient:      user.NewUserServiceClient(userConn),
    }
}

func (s *Server) Start() {
    router := gin.Default()

    userConn := mustConnect(s.cfg.UserServiceAddr)
    inventoryConn := mustConnect(s.cfg.InventoryServiceAddr)
    orderConn := mustConnect(s.cfg.OrderServiceAddr)

    userClient := pb.NewUserServiceClient(userConn)
    inventoryClient := pb.NewInventoryServiceClient(inventoryConn)
    orderClient := pb.NewOrderServiceClient(orderConn)

    h := handler.NewHandler(userClient, inventoryClient, orderClient)

    setupRoutes(router, h)

    log.Printf("API Gateway is running on port %s", s.cfg.HttpPort)
    if err := router.Run(s.cfg.HttpPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func mustConnect(addr string) *grpc.ClientConn {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to %s: %v", addr, err)
    }
    return conn
}

// setupRoutes настраивает маршруты для API Gateway
func setupRoutes(router *gin.Engine, h *handler.Handler) {
    v1 := router.Group("/api/v1")

	open := v1.Group("/")
    {
        open.POST("/register", h.RegisterUser)
        open.POST("/login", h.Login)
    }

    auth := v1.Group("/")
    auth.Use(middleware.AuthMiddleware())
    {
		// User routes
        v1.GET("/profile/:id", h.GetProfile)
        // Inventory routes
        v1.POST("/products", h.CreateProduct)
        v1.GET("/products/:id", h.GetProduct)
        v1.PATCH("/products/:id", h.UpdateProduct)
        v1.DELETE("/products/:id", h.DeleteProduct)
        v1.GET("/products", h.ListProducts)

        // Order routes
        v1.POST("/orders", h.CreateOrder)
        v1.GET("/orders/:id", h.GetOrder)
        v1.GET("/orders/user/:id", h.ListUserOrders)
        v1.PATCH("/orders/:id/status", h.UpdateOrderStatus)
    }
}