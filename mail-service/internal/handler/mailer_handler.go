package handler

import (
	"context"
	"smtp-service/internal/pb"
	"smtp-service/internal/usecase"
)

type MailHandler struct {
	pb.UnimplementedMailerServiceServer
	service usecase.SMTPServicer 
}

func NewMailHandler(s usecase.SMTPServicer) *MailHandler {
	return &MailHandler{service: s}
}

func (h *MailHandler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	err := h.service.SendEmail(req.ToEmail, req.Subject, req.HtmlBody)
	if err != nil {
		return &pb.SendEmailResponse{Status: "error", Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Status: "success", Message: "Email sent"}, nil
}