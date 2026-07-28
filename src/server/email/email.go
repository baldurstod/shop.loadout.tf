package email

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/mail"

	"gopkg.in/gomail.v2"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/release"
)

var dialer *gomail.Dialer
var from string
var host = "https://shop.loadout.tf"

func SetMailConfig(smtp config.SMTP) {
	dialer = gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	from = smtp.From

	if release.ReleaseMode != "true" {
		host = "https://shop.loadout.localhost:17830"
	}
}

func sendMail(from string, to string, subject string, contentType string, body string) error {
	// Check destination validity
	if _, err := mail.ParseAddress(to); err != nil {
		return err
	}

	if dialer == nil {
		return errors.New("dialer is nil. Did you forgot to init mail config ?")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody(contentType, body)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error while sending mail to %s: %w", to, err)
	}

	return nil
}

func SendMail(from string, to string, subject string, body string) error {
	return sendMail(from, to, subject, "text/plain", body)
}

func SendMailHtml(from string, to string, subject string, body string) error {
	return sendMail(from, to, subject, "text/html", body)
}

func SendMailVerification(to string, code string, text string) error {
	t, err := template.New("body").Parse(text)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]interface{}{
		"host":  host,
		"email": html.EscapeString(to),
		"code":  html.EscapeString(code),
	})
	if err != nil {
		return err
	}

	return SendMailHtml(from, to, "Loadout.tf: verify your email address", buf.String())
}
