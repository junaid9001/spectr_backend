package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// send email
func SendEmail(to, subject, body string) error {

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASS")

	if from == "" || password == "" {
		return fmt.Errorf("email env variables not set")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil

}
