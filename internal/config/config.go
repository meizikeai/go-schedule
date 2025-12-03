// internal/config/config.go
package config

import (
	"log"
	"os"

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
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
	PoolSize int      `mapstructure:"pool_size"`
}

type KafkaInstance struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func Load() *Config {
	var result *Config

	path := "."
	if os.Getenv("GO_ENV") == "debug" {
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

	return result
}
