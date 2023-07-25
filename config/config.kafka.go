package config

var kafkaConfig = map[string]string{
	"default-test":    "127.0.0.1:9092,127.0.0.1:9092,127.0.0.1:9092",
	"default-release": "127.0.0.1:9092,127.0.0.1:9092,127.0.0.1:9092",
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
