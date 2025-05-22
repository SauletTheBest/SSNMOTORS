package usecase

import (
	"fmt"
	"net/smtp"
)

type SMTPService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewSMTPService(host, port, username, password string) *SMTPService {
	return &SMTPService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *SMTPService) SendEmail(to, subject, html string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		html + "\r\n")

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	return smtp.SendMail(addr, auth, s.Username, []string{to}, msg)
}
