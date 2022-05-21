package email_sender

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"webserver/internal/config"
)

func send(receiverMail, subject, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", config.GetInstance().Email)
	m.SetHeader("To", receiverMail)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GetInstance().Email, config.GetInstance().EmailSecret)
	d.StartTLSPolicy = gomail.MandatoryStartTLS

	return d.DialAndSend(m)
}

func SendResetToken(receiverMail, token string) error {
	url := fmt.Sprintf("%s/password-reset?token=%s", config.GetInstance().Url, token)
	return send(receiverMail, "Password reset",
		fmt.Sprintf("Click here to reset your password: %s", url))
}

func SendVerificationToken(receiverMail, token string) error {
	url := fmt.Sprintf("%s/email-verification?token=%s", config.GetInstance().Url, token)
	return send(receiverMail, "Email verification",
		fmt.Sprintf("Click here to verify your email: %s", url))
}

func GenerateToken() string {
	length := 20
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
