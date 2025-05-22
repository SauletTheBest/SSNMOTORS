package handler

import (
	"context"
	"smtp-service/internal/pb"
	"smtp-service/internal/usecase"
)

type MailHandler struct {
	pb.UnimplementedMailerServiceServer
	smtpService *usecase.SMTPService
}

func NewMailHandler(service *usecase.SMTPService) *MailHandler {
	return &MailHandler{smtpService: service}
}

func (h *MailHandler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	err := h.smtpService.SendEmail(req.ToEmail, req.Subject, req.HtmlBody)
	if err != nil {
		return &pb.SendEmailResponse{Status: "error", Message: err.Error()}, nil
	}

	return &pb.SendEmailResponse{Status: "success", Message: "Email sent"}, nil
}
