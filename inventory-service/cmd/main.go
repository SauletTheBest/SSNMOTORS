package main

import (
    "fmt"
    "log"
    "net"
    "inventory-service/config"
    "inventory-service/internal/handler"
    "inventory-service/internal/pb"
    "inventory-service/internal/repository"
    "inventory-service/internal/usecase"
    "inventory-service/internal/queue"
    "github.com/nats-io/nats.go"
    "google.golang.org/grpc"
)

func main() {
    cfg := config.Load()
    defer cfg.Client.Disconnect(cfg.Ctx)

    coll := cfg.Client.Database(cfg.MongoDBName).Collection("products")

    repo := repository.NewMongoProductRepository(coll)
    uc   := usecase.NewProductUsecase(repo)
    h    := handler.NewProductHandler(uc)

    natsConn, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatalf("❌ NATS connection failed: %v", err)
    }
    defer natsConn.Close()

    consumer := queue.NewConsumer(natsConn, "order.created", uc)
    go func() {
        if err := consumer.Subscribe(cfg.Ctx); err != nil {
            log.Fatalf("❌ Failed to subscribe to order.created: %v", err)
        }
        log.Println("📥 NATS subscription active on 'order.created'")
    }()

    lis, err := net.Listen("tcp", ":"+cfg.Port)
    if err != nil {
        log.Fatalf("listen error: %v", err)
    }

    srv := grpc.NewServer()
    pb.RegisterInventoryServiceServer(srv, h)

    fmt.Printf("🔆 InventoryService on port %s\n", cfg.Port)
    if err := srv.Serve(lis); err != nil {
        log.Fatalf("serve error: %v", err)
    }
}
