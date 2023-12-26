package models

import (
	"fmt"

	"go-schedule/libs/tool"
	"go-schedule/libs/types"
)

var test = &types.MailMessage{
	From:    "from@163.com",
	To:      []string{"to@163.com"},
	Subject: "Welcome!",
	Data:    `<h3>Hello!</h3><p>This is a test mail message.</p>`,
}

func SendTestMail() {
	m := tool.GetMailClient("mail")
	message := tool.CreateMailMessage(test)

	err := m.Send(test.From, test.To, message)

	if err != nil {
		fmt.Printf("Mail sending failed, %v", err)
	} else {
		fmt.Println("Mail sent successfully")
	}
}
