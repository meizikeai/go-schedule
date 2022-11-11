package tool

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"go-schedule/libs/types"

	"github.com/go-redis/redis/v8"

	log "github.com/sirupsen/logrus"
)

var connRedis = types.OutConfRedis{
	MaxRetries:         3,
	DialTimeout:        5,
	ReadTimeout:        500,
	WriteTimeout:       500,
	PoolSize:           10,
	IdleTimeout:        300,
	IdleCheckFrequency: 60,
}
var fullDbRedis map[string][]*redis.Client
var redisConfig types.FullConfRedis

func HandleLocalRedisConfig() {
	var config types.FullConfRedis

	pwd, _ := os.Getwd()
	mode := GetMODE()

	address := strings.Join([]string{pwd, "/conf/", mode, ".redis.json"}, "")

	res, err := ioutil.ReadFile(address)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &config)

	if err != nil {
		log.Fatal(err)
	}

	redisConfig = config
}

func GetRedisClient(key string) *redis.Client {
	result := fullDbRedis[key]
	count := GetRandmod(len(result))

	return result[count]
}

func HandleRedisClient() {
	config := make(map[string][]*redis.Client)

	zookeeper := getZookeeperRedisConfig()
	local := getLocalRedisConfig()

	for k, v := range zookeeper {
		key := k + ".master"

		for _, addr := range v.Master {
			clients := handleRedisClient(addr, v.Password, v.Db)
			config[key] = append(config[key], clients)
		}
	}

	for k, v := range local {
		key := k + ".master"

		for _, addr := range v.Master {
			clients := handleRedisClient(addr, v.Password, v.Db)
			config[key] = append(config[key], clients)
		}
	}

	fullDbRedis = config
}

func getLocalRedisConfig() types.FullConfRedis {
	return redisConfig
}

func createRedisClient(config types.OutConfRedis) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:        time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolSize:           config.PoolSize * runtime.NumCPU(),
		PoolTimeout:        time.Duration(config.ReadTimeout)*time.Millisecond + time.Second,
		IdleTimeout:        time.Duration(config.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(config.IdleCheckFrequency) * time.Second,
	})

	_, err := db.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func handleRedisClient(addr string, password string, db int) *redis.Client {
	option := types.OutConfRedis{
		Addr:               addr,
		Password:           password,
		DB:                 db,
		MaxRetries:         connRedis.MaxRetries,
		PoolSize:           connRedis.PoolSize,
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

	Stdout("Redis Close")
}
