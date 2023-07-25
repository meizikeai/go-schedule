package config

import "go-schedule/libs/types"

var mysqlConfig = map[string]types.ConfMySQL{
	"default-test": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306"},
		Username: "test",
		Password: "test@123",
		Database: "test",
	},
	"default-release": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306", "127.0.0.1:3306"},
		Username: "test",
		Password: "test@123",
		Database: "test",
	},
}

func GetMySQLConfig() types.FullConfMySQL {
	result := types.FullConfMySQL{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = mysqlConfig[key]
	}

	return result
}
