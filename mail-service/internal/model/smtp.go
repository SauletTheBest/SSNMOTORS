package model

type MailerService interface {
	SendEmail(to, subject, body string) error
}
