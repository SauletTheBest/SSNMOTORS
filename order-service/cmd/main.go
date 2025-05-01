package main

import (
	"fmt"
	"log"
	"net"
	"order-service/config"
	"order-service/internal/handler"
	"order-service/internal/pb"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	client, err := mongo.Connect(cfg.Ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(cfg.Ctx)

	orderRepo := repository.NewMongoOrderRepository(client.Database(cfg.MongoDBName).Collection("orders"))
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	fmt.Println("🚀 OrderService running on port", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}