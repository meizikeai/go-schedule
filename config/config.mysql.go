package config

import "go-schedule/libs/types"

var mysqlConfig = map[string]types.ConfMySQL{
	"default-test": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306"},
		Username: "test",
		Password: "test@123",
		Database: "test",
	},
	"default-release": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306", "127.0.0.1:3306"},
		Username: "test",
		Password: "test@123",
		Database: "test",
	},
}

func GetMySQLConfig() map[string]types.ConfMySQL {
	result := map[string]types.ConfMySQL{}

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = mysqlConfig[key]
	}

	return result
}

var binlogConfig = map[string]types.ConfBinlog{
	"default-dev": {
		ServerId: 1,
		Flavor:   "mysql",
		Name:     "mysql-bin.000001",
		Pos:      4,
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "test@123",
	},
	"default-pro": {
		ServerId: 1,
		Flavor:   "mysql",
		Name:     "mysql-bin.000001",
		Pos:      4,
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "test@123",
	},
}

func GetBinlogConfig() map[string]types.ConfBinlog {
	result := make(map[string]types.ConfBinlog, 0)

	data := []string{
		"default",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = binlogConfig[key]
	}

	return result
}

// mysqldump -h [host] -P [port] -u [user] -p[password] [database] [table]
// mysqldump -h 10.99.248.94 -P 3306 -u root -p zkdash zd_snapshot
var canalConfig = map[string]types.ConfMySQLCanal{
	"zookeeper-dev": {
		ServerId: 582531,
		Flavor:   "mysql",
		Addr:     "10.99.248.94:3306",
		Username: "root",
		Password: "qcloudblueCITY40666",
	},
	"zookeeper-pro": {
		ServerId: 1,
		Flavor:   "mysql",
		Addr:     "127.0.0.1:3306",
		Username: "test",
		Password: "test@123",
	},
}

func GetCanalConfig() map[string]types.ConfMySQLCanal {
	result := make(map[string]types.ConfMySQLCanal, 0)

	data := []string{
		"zookeeper",
	}

	for _, v := range data {
		key := getKey(v)
		result[v] = canalConfig[key]
	}

	return result
}
