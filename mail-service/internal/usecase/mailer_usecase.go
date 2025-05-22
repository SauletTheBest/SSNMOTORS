package usecase

import (
	"fmt"
	"net/smtp"
)

// SMTPServicer определяет интерфейс для отправки email
type SMTPServicer interface {
	SendEmail(to, subject, body string) error
}

type SMTPService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewSMTPService(host, port, username, password string) *SMTPService {
	return &SMTPService{host, port, username, password}
}

func (s *SMTPService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.Username, []string{to}, msg)
}