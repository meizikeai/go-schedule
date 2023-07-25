package config

import "go-schedule/libs/types"

var elasticSearchConfig = map[string]types.ConfElasticSearch{
	"default-test": {
		Address:  []string{"127.0.0.1:9200"},
		Username: "test",
		Password: "test@123",
	},
	"default-release": {
		Address:  []string{"127.0.0.1:9200"},
		Username: "test",
		Password: "test@123",
	},
}

func GetElasticSearchConfig() types.FullConfElasticSearch {

	result := types.FullConfElasticSearch{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = elasticSearchConfig[key]
	}

	return result
}
