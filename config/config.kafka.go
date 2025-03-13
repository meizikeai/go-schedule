package config

var kafkaConfig = map[string]string{
	"default-dev": "10.99.0.3:9092",
	"default-pro": "127.0.0.1:9092",
}

var kafkaConsumerGroupConfig = map[string]map[string]string{
	"default-dev": {
		"assignor": "range",
		"groupID":  "text",
		"oldest":   "true",
		"version":  "3.9.0",
	},
	"default-pro": {
		"assignor": "range",
		"groupID":  "text",
		"oldest":   "true",
		"version":  "3.9.0",
	},
}

func GetKafkaConfig() map[string]string {
	result := map[string]string{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = kafkaConfig[key]
	}

	return result
}

func GetKafkaConsumerGroupConfig(key string) map[string]string {
	return kafkaConsumerGroupConfig[getKey(key)]
}
