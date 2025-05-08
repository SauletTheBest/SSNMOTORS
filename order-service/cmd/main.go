package main

import (
	"context"
	"log"
	"net"
	"order-service/config"
	"order-service/internal/handler"
	"order-service/internal/pb"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using default env")
	}
	cfg := config.Load()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection error: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("⚠️  MongoDB disconnect error: %v", err)
		}
	}()

	db := client.Database(cfg.MongoDB)
	orderCollection := db.Collection("orders")

	// Setup repo, usecase, handler
	orderRepo := repository.NewMongoOrderRepository(orderCollection)
	orderUC := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUC)

	// Create and register gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Start gRPC listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	go func() {
		log.Println("🚀 gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("👋 Goodbye!")
}
