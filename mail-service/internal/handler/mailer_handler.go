package handler

import (
    "context"
    "mail-service/internal/pb"
    "mail-service/internal/usecase"
    
)

type MailerHandler struct {
    pb.UnimplementedMailerServiceServer
    mailer *usecase.MailerSendService
}

func NewMailerHandler(mailer *usecase.MailerSendService) *MailerHandler {
    return &MailerHandler{mailer: mailer}
}

func (h *MailerHandler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
    err := h.mailer.SendEmail(
        req.FromName,
        req.FromEmail,
        req.ToEmail,
        req.Subject,
        req.HtmlBody,
    )
    if err != nil {
        return &pb.SendEmailResponse{
            Status:  "error",
            Message: err.Error(),
        }, nil
    }

    return &pb.SendEmailResponse{
        Status:  "success",
        Message: "Email sent successfully",
    }, nil
}
