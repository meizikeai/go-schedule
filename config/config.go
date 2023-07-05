package config

import (
	"os"

	"go-schedule/libs/types"
)

var zookeeperTest = []string{"127.0.0.1:2181"}
var zookeeperRelease = []string{"127.0.0.1:2181"}

var mysqlConfig = map[string]types.ConfMySQL{
	"default-test": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
	"default-release": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306", "127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
}

var redisConfig = map[string]types.ConfRedis{
	"default-test": {
		Master:   []string{"127.0.0.1:6379"},
		Password: "",
		Db:       0,
	},
	"default-release": {
		Master:   []string{"127.0.0.1:6379"},
		Password: "",
		Db:       0,
	},
}

var zookeeperConfig = types.ConfZookeeper{
	"mysql": {
		"test": "/kiss/love/test",
	},
	"redis": {
		"test": "/kiss/love/test",
	},
	"server": {
		"test": "/kiss/love/test",
	},
	"string": {
		"test": "/kiss/love/test",
	},
	"kafka": {
		"test": "/kiss/love/test",
	},
}

func getMode() string {
	mode := os.Getenv("GO_MODE")

	if mode == "" {
		mode = "test"
	}

	return mode
}

func isProduction() bool {
	result := false

	mode := os.Getenv("GO_MODE")

	if mode == "release" {
		result = true
	}

	return result
}

func GetMySQLConfig() types.FullConfMySQL {
	mode := getMode()
	result := types.FullConfMySQL{}

	data := []string{
		"default",
	}

	for _, v := range data {
		result[v] = mysqlConfig[v+"-"+mode]
	}

	return result
}

func GetRedisConfig() types.FullConfRedis {
	mode := getMode()
	result := types.FullConfRedis{}

	data := []string{
		"default",
	}

	for _, v := range data {
		result[v] = redisConfig[v+"-"+mode]
	}

	return result
}

func GetZookeeperConfig() types.ConfZookeeper {
	return zookeeperConfig
}

func GetZookeeperServers() []string {
	env := isProduction()

	result := zookeeperRelease

	if env == false {
		result = zookeeperTest
	}

	return result
}
