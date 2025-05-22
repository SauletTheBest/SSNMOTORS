package test

import (
	"context"
	"smtp-service/internal/handler"
	"smtp-service/internal/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSMTPService struct {
	mock.Mock
}

func (m *MockSMTPService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestSendEmailHandler(t *testing.T) {
	mockService := &MockSMTPService{}
	h := handler.NewMailHandler(mockService)

	// Настраиваем ожидания
	mockService.On("SendEmail", "test@example.com", "Test", "<b>Hello</b>").Return(nil)

	resp, err := h.SendEmail(context.Background(), &pb.SendEmailRequest{
		ToEmail:  "test@example.com",
		Subject:  "Test",
		HtmlBody: "<b>Hello</b>",
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Email sent", resp.Message)
	mockService.AssertExpectations(t)
}