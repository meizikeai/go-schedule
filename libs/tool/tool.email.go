package tool

import (
	"go-schedule/libs/types"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Client map[string]gomail.SendCloser
}

func NewEmail(data map[string]types.ConfMail) *Email {
	client := make(map[string]gomail.SendCloser)

	for k, v := range data {
		dialer := gomail.NewDialer(v.Host, v.Port, v.Username, v.Password)
		dialer.SSL = true

		d, err := dialer.Dial()

		if err != nil {
			panic(err)
		}

		client[k] = d
	}

	return &Email{
		Client: client,
	}
}

func (e *Email) Close() {
	for _, v := range e.Client {
		v.Close()
	}
}
