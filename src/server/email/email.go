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

	if release.ReleaseMode == "true" {
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

func SendMailVerification(to string, code string) error {
	var buf bytes.Buffer
	err := verifyTemplate.Execute(&buf, map[string]interface{}{
		"host":  host,
		"email": html.EscapeString(to),
		"code":  html.EscapeString(code),
	})
	if err != nil {
		return err
	}

	return SendMailHtml(from, to, "Loadout.tf: verify your email address", buf.String())
}

func createVerifyTemplate() *template.Template {
	t, err := template.New("body").Parse(verifyBody)
	if err != nil {
		panic(err)
	}
	return t
}

var verifyTemplate = createVerifyTemplate()

var verifyBody = `
	<html>
	<body>
	<h1>Verify your email address</h1>
	To finish setting up your account, we just need to make sure this email address is yours.<br>
	To verify your email address, <a href="{{.host}}/verify_email?code={{.code}}&email={{.email}}">click on this link.</a><br>
	If you didn't request this code, you can safely ignore this email.

	</body>
	</html>
	`
