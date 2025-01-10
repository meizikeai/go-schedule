package tool

import (
	"context"
	"runtime"
	"time"

	"go-schedule/libs/types"

	"github.com/go-redis/redis/v8"
)

type configRedis struct {
	Addr               string `json:"addr"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	DB                 int    `json:"db"`
	MaxRetries         int    `json:"max_retries"`
	PoolSize           int    `json:"pool_size"`
	MinIdleConns       int    `json:"min_idle_conns"`
	DialTimeout        int    `json:"dial_timeout"`
	ReadTimeout        int    `json:"read_timeout"`
	WriteTimeout       int    `json:"write_timeout"`
	IdleTimeout        int    `json:"idle_timeout"`
	IdleCheckFrequency int    `json:"idle_check_frequency"`
}

type Redis struct {
	Client map[string][]*redis.Client
}

var connRedis = configRedis{
	MaxRetries:         3,
	PoolSize:           20,
	MinIdleConns:       10,
	DialTimeout:        5,
	ReadTimeout:        500,
	WriteTimeout:       500,
	IdleTimeout:        300,
	IdleCheckFrequency: 60,
}

func NewRedisClient(data map[string]types.ConfRedis) *Redis {
	client := make(map[string][]*redis.Client)

	for k, v := range data {
		for _, addr := range v.Master {
			clients := handleRedisClient(addr, v.Password, v.Db)
			client[k] = append(client[k], clients)
		}
	}

	return &Redis{
		Client: client,
	}
}

func handleRedisClient(addr, password string, db int) *redis.Client {
	option := configRedis{
		Addr:               addr,
		Password:           password,
		DB:                 db,
		MaxRetries:         connRedis.MaxRetries,
		PoolSize:           connRedis.PoolSize * runtime.NumCPU(),
		MinIdleConns:       connRedis.MinIdleConns * runtime.NumCPU(),
		ReadTimeout:        connRedis.ReadTimeout,
		WriteTimeout:       connRedis.WriteTimeout,
		IdleTimeout:        connRedis.IdleTimeout,
		IdleCheckFrequency: connRedis.IdleCheckFrequency,
	}

	client := createRedisClient(option)

	return client
}

func createRedisClient(config configRedis) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		DialTimeout:        time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:        time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolTimeout:        time.Duration(config.ReadTimeout)*time.Millisecond + time.Second,
		IdleTimeout:        time.Duration(config.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(config.IdleCheckFrequency) * time.Second,
	})

	_, err := db.Ping(context.Background()).Result()

	if err != nil {
		panic(err)
	}

	return db
}

func (r *Redis) Close() {
	for _, val := range r.Client {
		for _, v := range val {
			v.Close()
		}
	}
}
