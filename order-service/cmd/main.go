package main

import (
	"fmt"
	"log"
	"net"

	"order-service/config"
	"order-service/internal/handler"
	"order-service/internal/pb"
	"order-service/internal/queue"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"order-service/internal/cache"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// MongoDB
	client, err := mongo.Connect(cfg.Ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(cfg.Ctx)

	// NATS
	nc, err := nats.Connect(cfg.NATSURL) 
	if err != nil {
		log.Fatalf("❌ NATS connection failed: %v", err)
	}
	defer nc.Close()

	publisher, err := queue.NewNATSPublisher(nc)
	if err != nil {
		log.Fatalf("❌ Failed to create NATS publisher: %v", err)
	}
	// Redis
	redisClient := cache.NewRedisClient()

	orderRepo := repository.NewMongoOrderRepository(client.Database(cfg.MongoDBName).Collection("orders"))
	orderUsecase := usecase.NewOrderUsecase(orderRepo, publisher , redisClient) 

	orderHandler := handler.NewOrderHandler(orderUsecase, publisher)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("❌ Listen error: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	fmt.Println("🚀 OrderService running on port", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ gRPC serve error: %v", err)
	}
}