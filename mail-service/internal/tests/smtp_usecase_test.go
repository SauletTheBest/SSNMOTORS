package test

import (
	"smtp-service/internal/usecase"
	"testing"
)

func TestSendEmail(t *testing.T) {
	s := usecase.NewSMTPService("smtp.gmail.com", "587", "your-email@gmail.com", "your-app-password")

	err := s.SendEmail("recipient@example.com", "Test Subject", "<h1>Hello</h1>")
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}
