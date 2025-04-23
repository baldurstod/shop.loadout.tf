package mail

import (
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"
	"shop.loadout.tf/src/server/config"
)

var dialer *gomail.Dialer

func SetMailConfig(smtp config.SMTP) {
	dialer = gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
}

func SendMail(from string, to string, subject string, body string) error {
	if dialer == nil {
		return errors.New("dialer is nil")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error while sending mail to %s: %w", to, err)
	}

	return nil
}
