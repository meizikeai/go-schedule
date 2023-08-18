package config

import "go-schedule/libs/types"

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

func GetRedisConfig() map[string]types.ConfRedis {
	result := map[string]types.ConfRedis{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = redisConfig[key]
	}

	return result
}
