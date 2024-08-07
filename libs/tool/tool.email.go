package tool

import (
	"go-schedule/config"
	"go-schedule/libs/types"

	"gopkg.in/gomail.v2"
)

var oneMail map[string]gomail.SendCloser

func (t *Tools) GetMailClient(key string) gomail.SendCloser {
	return oneMail[key]
}

func (t *Tools) HandleMailClient() {
	clients := make(map[string]gomail.SendCloser)

	local := config.GetMailConfig()

	for k, v := range local {
		client := gomail.NewDialer(v.Host, v.Port, v.Username, v.Password)
		client.SSL = true

		d, err := client.Dial()

		if err != nil {
			panic(err)
		}

		clients[k] = d
	}

	oneMail = clients

	t.Stdout("Mail Dialer is Connected")
}

func (t *Tools) CloseMail() {
	var err error

	for _, v := range oneMail {
		err = v.Close()

		if err != nil {
			break
		}
	}

	if err != nil {
		t.Stdout("Mail Dialer is Close")
	}
}

func (t *Tools) CreateMailMessage(e *types.MailMessage) *gomail.Message {
	m := gomail.NewMessage()

	if len(e.Cc) > 0 {
		m.SetHeader("Cc", e.Cc...)
	}

	for _, v := range e.File {
		m.Attach(v)
	}

	m.SetHeader("From", e.From)
	m.SetHeader("To", e.To...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody(
		"text/html",
		e.Data,
	)

	return m
}
