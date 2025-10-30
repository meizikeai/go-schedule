package repository

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type SendMail struct{}

func NewSendMail() *SendMail {
	return &SendMail{}
}

type mailMessage struct {
	From    string   `json:"from"`
	Subject string   `json:"subject"`
	Data    string   `json:"data"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	File    []string `json:"file"`
}

func (s *SendMail) SendTestMail() {
	m := tools.GetMailClient("mail")

	demo := mailMessage{
		From:    "from@163.com",
		To:      []string{"to@163.com"},
		Subject: "Welcome!",
		Data:    `<h3>Hello!</h3><p>This is a test mail message.</p>`,
	}
	message := s.CreateMailMessage(&demo)

	err := m.Send(demo.From, demo.To, message)

	if err != nil {
		fmt.Printf("Mail sending failed, %v", err)
	} else {
		fmt.Println("Mail sent successfully")
	}
}

func (s *SendMail) CreateMailMessage(e *mailMessage) *gomail.Message {
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
