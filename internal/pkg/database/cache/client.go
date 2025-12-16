// internal/pkg/database/cache/client.go
package cache

import (
	"math/rand"

	"go-schedule/internal/config"

	"github.com/go-redis/redis/v8"
)

type Clients struct {
	clients map[string][]*redis.Client
}

func NewClient(cfg *map[string][]config.RedisInstance) *Clients {
	c := &Clients{
		clients: make(map[string][]*redis.Client),
	}

	for key, value := range *cfg {
		for _, v := range value {
			db := createClient(&v)
			c.clients[key] = append(c.clients[key], db)
		}
	}

	return c
}

func createClient(option *config.RedisInstance) *redis.Client {
	cfg := &redis.Options{
		Addr:         option.Addrs[0],
		Password:     option.Password,
		DB:           option.DB,
		PoolSize:     option.PoolSize,
		MinIdleConns: option.MinIdleConns,
	}

	client := redis.NewClient(cfg)

	return client
}

func (c *Clients) Client(key string) *redis.Client {
	clients := c.clients[key]
	if len(clients) > 0 {
		return clients[rand.Intn(len(clients))]
	}

	return nil
}

func (c *Clients) Close() {
	for _, v := range c.clients {
		for _, c := range v {
			c.Close()
		}
	}
}
