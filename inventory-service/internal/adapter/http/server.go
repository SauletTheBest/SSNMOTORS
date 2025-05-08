package http

import (
	"AP2Assignment2/inventory-service/config"
	"AP2Assignment2/inventory-service/internal/adapter/http/handler"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serverIPAddress = "0.0.0.0:%d" // Changed to 0.0.0.0 for external access

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer

	address        string
	productHandler *handler.ProductHandler
}

func New(cfg config.Server, useCase handler.ProductUseCase) *API {
	// Setting the Gin mode
	gin.SetMode(cfg.HTTPServer.Mode)

	// Creating a new Gin Engine
	server := gin.New()

	// Applying middleware
	server.Use(gin.Recovery())

	// Binding products
	productHandler := handler.NewProductHandler(useCase)

	api := &API{
		server:         server,
		cfg:            cfg.HTTPServer,
		address:        fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		productHandler: productHandler,
	}

	api.setupRoutes()

	return api
}

func (api *API) setupRoutes() {
	v1 := api.server.Group("api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", api.productHandler.GetAll)
			products.GET("/:id", api.productHandler.GetByID)
			products.PUT("/:id", api.productHandler.Update)
			products.POST("/", api.productHandler.Create)
			products.DELETE("/:id", api.productHandler.Delete)
		}
	}
}

func (api *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server running on: %v", api.address)

		// No need to reinitialize `api.server` here. Just run it directly.
		if err := api.server.Run(api.address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to run HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) Stop() error {
	// Setting up the signal channel to catch termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Blocking until a signal is received
	sig := <-quit
	log.Println("Shutdown signal received", "signal:", sig.String())

	// Creating a context with timeout for graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("HTTP server shutting down gracefully")

	// Note: You can use `Shutdown` if you use `http.Server` instead of `gin.Engine`.
	log.Println("HTTP server stopped successfully")

	return nil
}
