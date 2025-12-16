// internal/config/common.go
package config

import (
	"log"
	"os"
	"slices"

	"github.com/spf13/viper"
)

type Config struct {
	App   App                        `mapstructure:"app"`
	Host  map[string]string          `mapstructure:"host"`
	Kafka map[string]KafkaInstance   `mapstructure:"kafka"`
	MySQL map[string][]MySQLInstance `mapstructure:"mysql"`
	Redis map[string][]RedisInstance `mapstructure:"redis"`
}

type App struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
}

type MySQLInstance struct {
	Master          []string `mapstructure:"master"`
	Slave           []string `mapstructure:"slave"`
	MaxIdleConns    int      `mapstructure:"max_idle_conns"`
	MaxOpenConns    int      `mapstructure:"max_open_conns"`
	ConnMaxLifetime int      `mapstructure:"conn_max_lifetime"`
}

type RedisInstance struct {
	Addrs        []string `mapstructure:"addrs"`
	Password     string   `mapstructure:"password"`
	DB           int      `mapstructure:"db"`
	PoolSize     int      `mapstructure:"pool_size"`
	MinIdleConns int      `mapstructure:"min_idle_conns"`
}

type KafkaInstance struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func Load() *Config {
	var result *Config
	var env = []string{"release", "test"}
	var mode = os.Getenv("GO_ENV")
	var path = "."

	if !slices.Contains(env, mode) {
		mode = "test"
	}

	if mode == "release" {
		path = path + "/release"
	} else {
		path = path + "/test"
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Read config failed: %v", err)
	}

	if err := viper.Unmarshal(&result); err != nil {
		log.Fatalf("Unmarshal config failed: %v", err)
	}

	result.App.Mode = mode

	return result
}
