package config

var mongodbConfig = map[string][]string{
	"default-test": {
		"mongodb://admin:123456@localhost:27017",
		"mongodb://admin:123456@localhost:27017",
		"mongodb://admin:123456@localhost:27017",
	},
	"default-release": {
		"mongodb://admin:123456@localhost:27017",
	},
}

func GetMongodbConfig() map[string][]string {
	result := map[string][]string{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = mongodbConfig[key]
	}

	return result
}
