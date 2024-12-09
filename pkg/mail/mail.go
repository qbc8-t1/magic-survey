package mail

import (
	"fmt"
	"net/smtp"
)

func SendMail(from string, pass string, to string, subject string, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil

}
