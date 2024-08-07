package tool

import (
	"time"

	"go-schedule/config"
	"go-schedule/libs/types"

	"github.com/go-zookeeper/zk"
)

type zookeeper struct {
	Client *zk.Conn
}

var (
	zookeeperService zookeeper
	zookeeperMySQL   map[string]types.ConfMySQL
	zookeeperRedis   map[string]types.ConfRedis
	zookeeperApi     map[string][]string
)

func (t *Tools) HandleZookeeperClient() {
	zk := t.getZookeeperService()
	option := config.ZookeeperConfig

	for key, val := range option {
		if key == "mysql" {
			config := make(map[string]types.ConfMySQL, 0)

			for k, v := range val {
				back := zk.GetZookeeperMySQLConfig(v)
				config[k] = back
			}

			zookeeperMySQL = config
		} else if key == "redis" {
			config := make(map[string]types.ConfRedis, 0)

			for k, v := range val {
				back := zk.GetZookeeperRedisConfig(v)
				config[k] = back
			}

			zookeeperRedis = config
		} else if key == "api" {
			config := make(map[string][]string, 0)

			for k, v := range val {
				back := zk.GetZookeeperApiConfig(v)
				config[k] = back
			}

			zookeeperApi = config
		}
	}

	defer zk.Close()
}

func (t *Tools) getMySQLConfig() map[string]types.ConfMySQL {
	return zookeeperMySQL
}

func (t *Tools) getRedisConfig() map[string]types.ConfRedis {
	return zookeeperRedis
}

func (t *Tools) newZookeeper(servers []string) *zk.Conn {
	client, _, err := zk.Connect(servers, 4*time.Second)

	if err != nil {
		panic(err)
	}

	return client
}

func (t *Tools) getZookeeperService() *zookeeper {
	config := config.GetZookeeperConfig()
	zookeeperService.Client = t.newZookeeper(config)

	return &zookeeperService
}

func (z *zookeeper) Get(path string) string {
	v, _, err := z.Client.Get(path)

	if err != nil {
		panic(err)
	}

	return string(v)
}

func (z *zookeeper) Children(path string) []string {
	res, _, err := z.Client.Children(path)

	if err != nil {
		panic(err)
	}

	return res
}

func (z *zookeeper) filterData(path string) []string {
	result := make([]string, 0)

	data := z.Children(path)

	for _, v := range data {
		key := path + "/" + v

		back := z.Get(key)

		if back == "0" {
			result = append(result, v)
		}
	}

	return result
}

func (z *zookeeper) Close() {
	z.Client.Close()
}

func (z *zookeeper) GetZookeeperApiConfig(path string) []string {
	return z.filterData(path)
}

func (z *zookeeper) GetZookeeperRedisConfig(path string) types.ConfRedis {
	redis := types.ConfRedis{
		Master: z.filterData(path),
	}

	return redis
}

func (z *zookeeper) GetZookeeperMySQLConfig(path string) types.ConfMySQL {
	mysql := types.ConfMySQL{
		Master:   nil,
		Slave:    nil,
		Username: "",
		Password: "",
		Database: "",
	}

	back := z.Children(path)

	for _, val := range back {
		key := path + "/" + val

		switch val {
		case "master":
			mysql.Master = z.filterData(key)
		case "slave":
			mysql.Slave = z.filterData(key)
		case "username":
			mysql.Username = z.Get(key)
		case "password":
			mysql.Password = z.Get(key)
		case "database":
			mysql.Database = z.Get(key)
		}
	}

	return mysql
}
