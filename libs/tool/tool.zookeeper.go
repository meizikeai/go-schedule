package tool

import (
	"time"

	"go-schedule/libs/types"

	"github.com/go-zookeeper/zk"
)

type zookeeper struct {
	Client *zk.Conn
}

var zookeeperService zookeeper

func init() {
	// // not use
	// config := config.GetZookeeperConfig()
	// zookeeperService.Client = newZookeeper(config)
}

func newZookeeper(servers []string) *zk.Conn {
	client, _, err := zk.Connect(servers, 4*time.Second)

	if err != nil {
		panic(err)
	}

	return client
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

func (z *zookeeper) Close() {
	z.Client.Close()
}

// demo - get redis config
func (z *zookeeper) GetZookeeperRedisConfig(name, path string) map[string]types.ConfRedis {
	redis := types.ConfRedis{}

	back := z.Children(path)

	redis.Master = back

	result := map[string]types.ConfRedis{
		name: redis,
	}

	defer z.Close()

	return result
}
