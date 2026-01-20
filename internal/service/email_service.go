package service

import (
	"fmt"
	"os"
	"strconv"
)

type EmailService struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService() *EmailService {
	port := 587
	portStr := os.Getenv("EMAIL_PORT")
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &EmailService{
		From:     os.Getenv("EMAIL_FROM"),
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}

func (s *EmailService) SendMail(to, subject, body string) error {
	if s.Host == "" || s.Username == "" {
		// If email config is not set, skip sending (useful for dev)
		fmt.Printf("Email would be sent to %s with subject: %s\n", to, subject)
		return nil
	}

	// TODO: Implement actual email sending with gomail when package is added
	// For now, just log that email would be sent
	fmt.Printf("Sending email to %s with subject: %s\n", to, subject)
	return nil
}
