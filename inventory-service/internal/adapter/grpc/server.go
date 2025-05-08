package grpc

import (
	"AP2Assignment2/inventory-service/config"
	"AP2Assignment2/inventory-service/internal/usecase"
	proto "AP2Assignment2/inventory-service/protos/gen"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serverIPAddress = "0.0.0.0:%d"

type ServerAPI struct {
	grpcServer       *grpc.Server
	cfg              config.GRPCServer
	address          string
	inventoryHandler *InventoryGRPCServer
}

func New(cfg config.Server, productUsecase *usecase.Product) *ServerAPI {
	grpcServer := grpc.NewServer()

	inventoryHandler := NewInventoryGRPCServer(productUsecase)
	proto.RegisterInventoryServiceServer(grpcServer, inventoryHandler)

	server := &ServerAPI{
		grpcServer:       grpcServer,
		cfg:              cfg.GRPCServer,
		address:          fmt.Sprintf(serverIPAddress, cfg.GRPCServer.Port),
		inventoryHandler: inventoryHandler,
	}

	return server
}

func (s *ServerAPI) Run(errCh chan<- error) {
	go func() {
		log.Printf("gRPC server running on: %v", s.address)

		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			errCh <- fmt.Errorf("failed to listen on %s: %w", s.address, err)
			return
		}

		if err := s.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- fmt.Errorf("failed to run gRPC server: %w", err)
			return
		}
	}()
}

func (s *ServerAPI) Stop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("Shutdown signal received", "signal:", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("gRPC server shutting down gracefully")

	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("gRPC server stopped successfully")
	case <-ctx.Done():
		log.Println("gRPC server shutdown timed out, forcing stop")
		s.grpcServer.Stop()
	}

	return nil
}
