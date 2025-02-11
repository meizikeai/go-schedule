package config

var zookeeperConfig = map[string][]string{
	"zookeeper-dev": {"127.0.0.1:2181"},
	"zookeeper-pro": {"127.0.0.1:2181", "127.0.0.1:2181", "127.0.0.1:2181"},
}

var ZookeeperConfig = map[string]map[string]string{
	"mysql": {
		"tests": "/blue/backend/umem/tests",
	},
	"redis": {
		"users": "/blue/backend/umem/users",
	},
}

func GetZookeeperConfig(key string) []string {
	return zookeeperConfig[getKey(key)]
}
