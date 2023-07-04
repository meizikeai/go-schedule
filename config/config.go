package config

import (
	"os"

	"go-schedule/libs/types"
)

var zookeeperTest = []string{"127.0.0.1:2181"}
var zookeeperRelease = []string{"127.0.0.1:2181"}

var mysqlConfig = map[string]types.ConfMySQL{
	"ailab-test": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
	"ailab-release": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306", "127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
}

var redisConfig = map[string]types.ConfRedis{
	"ailab-test": {
		Master:   []string{"127.0.0.1:6379"},
		Password: "",
		Db:       0,
	},
	"ailab-release": {
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

func isProduction() bool {
	result := false

	mode := os.Getenv("ALP_MODE")

	if mode == "release" {
		result = true
	}

	return result
}

func GetMySQLConfig() types.FullConfMySQL {
	env := isProduction()

	confMySQL := mysqlConfig["ailab-release"]

	if env == false {
		confMySQL = mysqlConfig["ailab-test"]
	}

	result := types.FullConfMySQL{
		"default": confMySQL,
	}

	return result
}

func GetRedisConfig() types.FullConfRedis {
	env := isProduction()

	confRedis := redisConfig["ailab-release"]

	if env == false {
		confRedis = redisConfig["ailab-test"]
	}

	result := types.FullConfRedis{
		"default": confRedis,
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
