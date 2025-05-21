package main

import (
	"log"
	"mail-service/config"
	"mail-service/internal/handler"
	"mail-service/internal/pb"
	"mail-service/internal/usecase"
	"net"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	mailer := usecase.NewMailerSendService(cfg.MailerSendAPIKey)

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	mailerHandler := handler.NewMailerHandler(mailer)
	pb.RegisterMailerServiceServer(s, mailerHandler)

	log.Println("✅ gRPC Mailer Service running on port", cfg.ServerPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
