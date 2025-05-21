package usecase

import (
    "github.com/go-resty/resty/v2"
    "log"
    "fmt"

)

type MailerSendService struct {
    ApiKey string
}

func NewMailerSendService(apiKey string) *MailerSendService {
    return &MailerSendService{ApiKey: apiKey}
}

func (m *MailerSendService) SendEmail(fromName, fromEmail, toEmail, subject, htmlBody string) error {
    client := resty.New()

    payload := map[string]interface{}{
        "from": map[string]string{
            "email": fromEmail,
            "name":  fromName,
        },
        "to": []map[string]string{
            {
                "email": toEmail,
            },
        },
        "subject": subject,
        "html":    htmlBody,
    }

    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+m.ApiKey).
        SetHeader("Content-Type", "application/json").
        SetBody(payload).
        Post("https://api.mailersend.com/v1/email")

    if err != nil {
        return err
    }

    if resp.IsError() {
    log.Println("MailerSend Error:", string(resp.Body()))
    return fmt.Errorf("MailerSend API error: %s", string(resp.Body()))
}

    log.Println("📧 Email sent to:", toEmail)
    return nil
}
