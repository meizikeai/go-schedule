package tool

import (
	"time"

	"github.com/go-zookeeper/zk"
)

type Zookeeper struct {
	Client *zk.Conn
}

func NewZookeeper(servers []string) *Zookeeper {
	result, _, err := zk.Connect(servers, 4*time.Second)

	if err != nil {
		panic(err)
	}

	return &Zookeeper{
		Client: result,
	}
}

func (z *Zookeeper) Get(path string) string {
	v, _, err := z.Client.Get(path)

	if err != nil {
		panic(err)
	}

	return string(v)
}

func (z *Zookeeper) Children(path string) []string {
	res, _, err := z.Client.Children(path)

	if err != nil {
		panic(err)
	}

	return res
}

func (z *Zookeeper) Close() {
	z.Client.Close()
}
