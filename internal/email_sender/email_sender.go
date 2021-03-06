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

func SendDenial(receiverMail, taskName, adminName, adminEmail, reason string) error {
	return send(receiverMail, "Your task has been moved to unpublished",
		fmt.Sprintf("Your task %q has been moved to the unpublished state by %q(%s). The reason for it is: %q",
			taskName, adminName, adminEmail, reason))
}

func SendOnUserSolutionDeletion(receiverMail, code, adminName, adminEmail, reason string) error {
	return send(receiverMail, "Your solution has been removed from the statistic",
		fmt.Sprintf("Your solution has been removed from the statistic by %q(%s). The reason for it is: %q.\n"+
			"Code of your solution:\n\n%s", adminName, adminEmail, reason, code))
}

func SendOnTaskApproval(receiverMail, adminName, adminEmail, taskName string) error {
	return send(receiverMail, "Your task has been approved",
		fmt.Sprintf("Your task %q has been approved by %q(%s)", taskName, adminName, adminEmail))
}

func GenerateToken() string {
	length := 20
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
