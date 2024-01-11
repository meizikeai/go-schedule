package config

var zookeeperTest = []string{"127.0.0.1:2181"}
var zookeeperRelease = []string{"127.0.0.1:2181"}

var ZookeeperConfig = map[string]map[string]string{
	"mysql": {
		"relation": "/blued/backend/umem/relation",
	},
	"redis": {
		"users": "/blued/backend/umem/users",
	},
}

func GetZookeeperConfig() []string {
	env := isProduction()

	result := zookeeperRelease

	if env == false {
		result = zookeeperTest
	}

	return result
}
