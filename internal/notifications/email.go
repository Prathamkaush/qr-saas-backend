package notifications

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailSender struct {
	From     string
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailSender(from, host, port, username, password string) *EmailSender {
	return &EmailSender{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *EmailSender) Send(to string, subject string, body string) error {
	addr := s.Host + ":" + s.Port

	msg := "From: " + s.From + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	err := smtp.SendMail(addr, auth, s.From, []string{to}, []byte(msg))
	if err != nil {
		log.Println("❌ Email send failed:", err)
		return err
	}

	fmt.Println("📨 Email sent to:", to)
	return nil
}
