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

    "google.golang.org/grpc"
)

func main() {
    cfg := config.Load()
    defer cfg.Client.Disconnect(cfg.Ctx)

    coll := cfg.Client.Database(cfg.MongoDBName).Collection("products")

    repo := repository.NewMongoProductRepository(coll)
    uc   := usecase.NewProductUsecase(repo)
    h    := handler.NewProductHandler(uc)

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
