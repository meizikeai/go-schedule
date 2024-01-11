package tool

import (
	"context"
	"runtime"
	"time"

	"go-schedule/config"
	"go-schedule/libs/types"

	"github.com/go-redis/redis/v8"
)

var connRedis = types.OutConfRedis{
	MaxRetries:         3,
	PoolSize:           20,
	MinIdleConns:       10,
	DialTimeout:        5,
	ReadTimeout:        500,
	WriteTimeout:       500,
	IdleTimeout:        300,
	IdleCheckFrequency: 60,
}
var fullDbRedis map[string][]*redis.Client

func GetRedisClient(key string) *redis.Client {
	result := fullDbRedis[key]
	count := GetRandmod(len(result))

	return result[count]
}

func HandleRedisClient() {
	client := make(map[string][]*redis.Client)

	// local := getRedisConfig()
	local := config.GetRedisConfig()

	for k, v := range local {
		for _, addr := range v.Master {
			clients := handleRedisClient(addr, v.Password, v.Db)
			client[k] = append(client[k], clients)
		}
	}

	fullDbRedis = client

	Stdout("Redis is Connected")
}

func createRedisClient(config types.OutConfRedis) *redis.Client {
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

func handleRedisClient(addr, password string, db int) *redis.Client {
	option := types.OutConfRedis{
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

func CloseRedis() {
	for _, val := range fullDbRedis {
		for _, v := range val {
			v.Close()
		}
	}

	Stdout("Redis is Close")
}
