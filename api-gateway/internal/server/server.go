package server

import (
	"fmt"
	"api-gateway/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Server encapsulates the Gin router
type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

// NewServer initializes the server with routes and middleware
func NewServer(cfg *config.Config) *Server {
	r := gin.New()

	// Initialize handlers
	h, err := handler.NewHandler(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize handler: %v", err))
	}

	// Apply global middleware
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.TelemetryMiddleware())

	// Define routes
	api := r.Group("/api")

	// Unprotected routes (no authentication required)
	api.POST("/users/register", h.RegisterUser)
	api.POST("/users/authenticate", h.AuthenticateUser)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Inventory routes
		protected.POST("/inventory", h.CreateProduct)
		protected.GET("/inventory/:id", h.GetProduct)
		protected.PUT("/inventory/:id", h.UpdateProduct)
		protected.DELETE("/inventory/:id", h.DeleteProduct)
		protected.GET("/inventory", h.ListProducts)

		// Order routes	
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)
		protected.PUT("/orders/:id/status", h.UpdateOrderStatus)
		protected.GET("/orders", h.ListUserOrders)

		// User routes
		protected.GET("/users/:id", h.GetUserProfile)
		// Mailer routes
		protected.POST("/email", h.SendEmail)
	}

	return &Server{
		router: r,
		cfg:    cfg,
	}
}

// Start runs the server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.router.Run(addr)
}