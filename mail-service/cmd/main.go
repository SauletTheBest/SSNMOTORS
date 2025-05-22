package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"smtp-service/internal/handler"
	"smtp-service/internal/pb"
	"smtp-service/internal/usecase"
	"smtp-service/config"
)

func main() {

	cfg := config.LoadConfig()

	// Создаём SMTP-сервис
	smtpService := usecase.NewSMTPService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	mailHandler := handler.NewMailHandler(smtpService)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMailerServiceServer(grpcServer, mailHandler)

	log.Println("SMTP gRPC server running on port 50054...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
