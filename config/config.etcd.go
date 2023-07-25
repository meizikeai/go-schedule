package config

import "go-schedule/libs/types"

var etcdConfig = map[string]types.ConfEtcd{
	"default-test": {
		Address:  []string{"127.0.0.1:2379"},
		Username: "test",
		Password: "test@123",
	},
	"default-release": {
		Address:  []string{"127.0.0.1:2379"},
		Username: "test",
		Password: "test@123",
	},
}

func GetEtcdConfig() types.ConfEtcd {
	result := types.ConfEtcd{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result = etcdConfig[key]
	}

	return result
}
