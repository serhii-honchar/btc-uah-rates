package email

import (
	"btc-uah-rates/pkg/utils"
	"log"
	"net/smtp"
	"strings"
)

type EmailService interface {
	SendEmails(subject, body string, recipients []string) error
}

type emailServiceImpl struct {
	smtpServer string
	smtpPort   string
	sender     string
	password   string
}

func NewEmailService() EmailService {
	host := utils.GetEnvOrDefault("EMAIL_HOST", "localhost")
	port := utils.GetEnvOrDefault("EMAIL_PORT", "1025")
	username := utils.GetEnvOrDefault("EMAIL_USERNAME", "AB")
	password := utils.GetEnvOrDefault("EMAIL_PASSWORD", "password")

	return &emailServiceImpl{
		smtpServer: host,
		smtpPort:   port,
		sender:     username,
		password:   password,
	}
}

func (e *emailServiceImpl) SendEmails(subject, body string, recipients []string) error {
	to := strings.Join(recipients, ",")

	message := "From: " + e.sender + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"\n" +
		body

	err := smtp.SendMail(e.smtpServer+":"+e.smtpPort, nil, e.sender, recipients, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
