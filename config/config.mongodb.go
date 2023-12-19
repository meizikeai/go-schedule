package config

import "go-schedule/libs/types"

var mongodbConfig = map[string]types.ConfMongoDB{
	"mongo-test": {
		Master: "mongodb://user:123456@127.0.0.1:27017/test?replicaSet=xxx&authSource=admin",
		Slave:  "mongodb://user:123456@127.0.0.1:27017/test?replicaSet=xxx&authSource=admin&readPreference=secondaryPreferred",
	},
	"mongo-release": {
		Master: "mongodb://user:123456@127.0.0.1:27017/test?replicaSet=xxx&authSource=admin",
		Slave:  "mongodb://user:123456@127.0.0.1:27017/test?replicaSet=xxx&authSource=admin&readPreference=secondaryPreferred",
	},
}

func GetMongodbConfig() map[string]types.ConfMongoDB {
	result := map[string]types.ConfMongoDB{}

	data := []string{
		"mongo",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = mongodbConfig[key]
	}

	return result
}
