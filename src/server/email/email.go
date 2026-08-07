package email

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/mail"

	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"
	"gopkg.in/gomail.v2"
	assets "shop.loadout.tf"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/release"
)

var dialer *gomail.Dialer
var from string
var to string
var host = "https://shop.loadout.tf"

func GetMailOrigin() string {
	return from
}
func GetMailDestination() string {
	return to
}

func SetMailConfig(smtp config.SMTP) {
	dialer = gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	from = smtp.From
	to = smtp.To

	if release.ReleaseMode != "true" {
		host = "https://shop.loadout.localhost:17830"
	}
}

func sendMail(from string, to string, subject string, contentType string, body string, process func(m *gomail.Message)) error {
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
	if process != nil {
		process(message)
	}

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error while sending mail to %s: %w", to, err)
	}

	return nil
}

func SendMail(from string, to string, subject string, body string, process func(m *gomail.Message)) error {
	return sendMail(from, to, subject, "text/plain", body, process)
}

func SendMailHtml(from string, to string, subject string, body string, process func(m *gomail.Message)) error {
	return sendMail(from, to, subject, "text/html", body, process)
}

func SendMailVerification(to string, code string, text string) error {
	t, err := template.New("body").Parse(text)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]any{
		"host":  host,
		"email": html.EscapeString(to),
		"code":  html.EscapeString(code),
	})
	if err != nil {
		return err
	}

	return SendMailHtml(from, to, "Loadout.tf: verify your email address", buf.String(), nil)
}

func SendOrderMail(user model.User, order model.Order) error {
	funcMap := template.FuncMap{
		"formatprice": func(retailPrice decimal.Decimal, c string) string {
			locale := currency.NewLocale(user.Locale)
			formatter := currency.NewFormatter(locale)
			amount, _ := currency.NewAmount(retailPrice.String(), order.Currency)
			return formatter.Format(amount)
		},
		"itemtotal": func(item model.OrderItem) decimal.Decimal {
			return decimal.NewFromInt(int64(item.Quantity)).Mul(item.GetRetailPrice())
		},
		"itemurl": func(item model.OrderItem) string {
			return host + "/@product/" + item.ProductID
		},
		"orderurl": func(order model.Order) string {
			return host + "/@order/" + order.ID
		},
	}

	t, err := template.New("order.html").Funcs(funcMap).ParseFS(&assets.TemplateAssets, "src/templates/order.html")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	var data = map[string]any{
		"name":     user.DisplayName,
		"items":    order.Items,
		"currency": order.Currency,
		"order":    order,
	}

	if order.SameBillingAddress {
		data["addresses"] = []model.Address{order.ShippingAddress, order.ShippingAddress}
	} else {
		data["addresses"] = []model.Address{order.BillingAddress, order.ShippingAddress}
	}

	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	log.Println(buf.String())

	embed := func(m *gomail.Message) {
		embed(m, "src/templates/order.css", "order_style")
	}

	return SendMailHtml(from, to, "Thank you for your purchase on shop.loadout.tf", buf.String(), embed)
}

func embed(message *gomail.Message, path string, cid string) error {
	b, err := assets.TemplateAssets.ReadFile(path)
	if err != nil {
		return err
	}

	message.Embed(path,
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(b)
			return err
		}),
		gomail.Rename(cid),
		gomail.SetHeader(map[string][]string{
			"Content-Type": {"text/css"},
		}),
	)
	return nil
}
