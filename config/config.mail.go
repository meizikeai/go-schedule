package config

import "go-schedule/libs/types"

var mailConfig = map[string]types.ConfMail{
	"mail-test": {
		Host:     "smtp.example.com",
		Port:     465,
		Username: "your@example.com",
		Password: "123456",
	},
	"mail-release": {
		Host:     "smtp.example.com",
		Port:     465,
		Username: "your@example.com",
		Password: "123456",
	},
}

func GetMailConfig() map[string]types.ConfMail {
	result := map[string]types.ConfMail{}

	data := []string{
		"mail",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = mailConfig[key]
	}

	return result
}
