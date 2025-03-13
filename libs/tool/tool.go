package tool

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"time"

	"go-schedule/config"
	"go-schedule/libs/types"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
)

var (
	dbMySQLCache            map[string][]*sql.DB
	dbMySQLReplicationCache map[string]*replication.BinlogSyncer
	dbMySQLCanalCache       map[string]*canal.Canal
	dbRedisCache            map[string][]*redis.Client
	dbMongoDBCache          map[string]*mongo.Client
	zookeeperMySQL          map[string]types.ConfMySQL
	zookeeperRedis          map[string]types.ConfRedis
	zookeeperApi            map[string][]string
	zookeeperConfig         map[string]string
	esClient                map[string][]*elasticsearch.Client
	kafkaProducer           map[string]sarama.AsyncProducer
	kafkaConsumer           map[string]sarama.Consumer
	kafkaConsumerGroup      map[string]sarama.ConsumerGroup
	emailClient             map[string]gomail.SendCloser
)

type Tools struct{}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func (t *Tools) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetAppName() string {
	return os.Getenv("CS_APP")
}

func GetGoEnv() string {
	return os.Getenv("GO_ENV")
}

// crontab
func (t *Tools) HandleCron(time string, fn func()) *cron.Cron {
	result := cron.New()

	_, err := result.AddFunc(time, fn)

	if err != nil {
		fmt.Println(err)
	}

	result.Start()

	return result
}

// mysql
func (t *Tools) GetMySQLClient(key string) *sql.DB {
	result := dbMySQLCache[key]
	index := t.GetRandmod(len(result))

	return result[index]
}

func (t *Tools) HandleMySQLClient() {
	config := zookeeperMySQL
	result := NewMySQLClient(config)

	dbMySQLCache = result.Client

	t.Stdout("MySQL is Connected")
}

func (t *Tools) CloseMySQL() {
	for _, val := range dbMySQLCache {
		for _, v := range val {
			v.Close()
		}
	}

	t.Stdout("MySQL is Close")
}

// mysql binlog - replication
func (t *Tools) GetReplicationMySQLClient(key string) *replication.BinlogSyncer {
	return dbMySQLReplicationCache[key]
}

func (t *Tools) HandleReplicationMySQLClient() {
	config := config.GetBinlogConfig()
	result := NewReplicationMySQLClient(config)

	dbMySQLReplicationCache = result.Client

	t.Stdout("MySQL Binlog Replication is Connected")
}

// mysql binlog - canal
func (t *Tools) GetCanalMySQLClient(key string) *canal.Canal {
	return dbMySQLCanalCache[key]
}

func (t *Tools) HandleCanalMySQLClient() {
	config := config.GetCanalConfig()
	result := NewCanalMySQLClient(config)

	dbMySQLCanalCache = result.Client

	t.Stdout("MySQL Binlog Canal is Connected")
}

// redis
func (t *Tools) GetRedisClient(key string) *redis.Client {
	result := dbRedisCache[key]
	index := t.GetRandmod(len(result))

	return result[index]
}

func (t *Tools) HandleRedisClient() {
	config := zookeeperRedis
	result := NewRedisClient(config)

	dbRedisCache = result.Client

	t.Stdout("Redis is Connected")
}

func (t *Tools) CloseRedis() {
	for _, val := range dbRedisCache {
		for _, v := range val {
			v.Close()
		}
	}

	t.Stdout("Redis is Close")
}

// api
func (t *Tools) GetApiClient(key string) string {
	result := ""

	for k, v := range zookeeperApi {
		if k == key {
			i := t.GetRandmod(len(v))
			result = v[i]
		}
	}

	return result
}

// config
func (t *Tools) GetConfigData(key string) string {
	return zookeeperConfig[key]
}

// zookeeper
func (t *Tools) HandleZookeeperClient() {
	servers := config.GetZookeeperConfig("zookeeper")
	zookeeper := NewZookeeper(servers)

	for key, val := range config.ZookeeperConfig {
		if key == "mysql" {
			config := make(map[string]types.ConfMySQL, 0)

			for k, v := range val {
				mysql := types.ConfMySQL{
					Master:   nil,
					Slave:    nil,
					Username: "",
					Password: "",
					Database: "",
				}

				back := zookeeper.Children(v)

				for _, val := range back {
					key := v + "/" + val

					switch val {
					case "master":
						mysql.Master = filterData(zookeeper, key)
					case "slave":
						mysql.Slave = filterData(zookeeper, key)
					case "username":
						mysql.Username = zookeeper.Get(key)
					case "password":
						mysql.Password = zookeeper.Get(key)
					case "database":
						mysql.Database = zookeeper.Get(key)
					}
				}

				config[k] = mysql
			}

			zookeeperMySQL = config
		} else if key == "redis" {
			config := make(map[string]types.ConfRedis, 0)

			for k, v := range val {
				redis := types.ConfRedis{
					Master: filterData(zookeeper, v),
				}

				config[k] = redis
			}

			zookeeperRedis = config
		} else if key == "api" {
			config := make(map[string][]string, 0)

			for k, v := range val {
				back := filterData(zookeeper, v)
				config[k] = back
			}

			zookeeperApi = config
		} else if key == "config" {
			config := make(map[string]string, 0)

			for k, v := range val {
				back := zookeeper.Get(v)
				config[k] = back
			}

			zookeeperConfig = config
		}
	}

	defer zookeeper.Close()
}

func filterData(z *Zookeeper, path string) []string {
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

// es
func (t *Tools) GetElasticSearchClient(key string) *elasticsearch.Client {
	result := esClient[key]
	count := t.GetRandmod(len(result))

	return result[count]
}

func (t *Tools) HandleElasticSearchClient() {
	local := config.GetElasticSearchConfig()
	result := NewElasticSearch(local)

	esClient = result.Client

	t.Stdout("ElasticSearch is Connected")
}

// kafka
// kafka producer
func (t *Tools) HandleKafkaProducerClient() {
	config := config.GetKafkaConfig()
	result := NewKafkaProducer(config)

	kafkaProducer = result.Client

	t.Stdout("Kafka Producer is Connected")
}

func (t *Tools) GetKafkaProducerClient(broker string) sarama.AsyncProducer {
	if kafkaProducer[broker] == nil {
		panic("Kafka Producer is not connected")
	}

	return kafkaProducer[broker]
}

// kafka consumer
func (t *Tools) HandleKafkaConsumerClient() {
	config := config.GetKafkaConfig()
	result := NewKafkaConsumer(config)

	kafkaConsumer = result.Client

	t.Stdout("Kafka Consumer is Connected")
}

func (t *Tools) GetKafkaConsumerClient(broker string) sarama.Consumer {
	if kafkaConsumer[broker] == nil {
		panic("Kafka Consumer is not connected")
	}

	return kafkaConsumer[broker]
}

// kafka consumer group
func (t *Tools) HandleKafkaConsumerGroupClient() {
	config := config.GetKafkaConfig()
	result := NewKafkaConsumerGroup(config)

	kafkaConsumerGroup = result.Client

	t.Stdout("Kafka Consumer Group is Connected")
}

func (t *Tools) GetKafkaConsumerGroupClient(broker string) sarama.ConsumerGroup {
	if kafkaConsumerGroup[broker] == nil {
		panic("Kafka Consumer Group  is not connected")
	}

	return kafkaConsumerGroup[broker]
}

func (t *Tools) CloseKafka() {
	if len(kafkaProducer) > 0 {
		for _, v := range kafkaProducer {
			v.Close()
		}
		t.Stdout("Kafka Producer is Close")
	}

	if len(kafkaConsumer) > 0 {
		for _, v := range kafkaConsumer {
			v.Close()
		}
		t.Stdout("Kafka Consumer is Close")
	}

	if len(kafkaConsumerGroup) > 0 {
		for _, v := range kafkaConsumerGroup {
			v.Close()
		}
		t.Stdout("Kafka Consumer Group is Close")
	}
}

// mongodb
func (t *Tools) GetMongoDBClient(key, database, table string) *mongo.Collection {
	client := dbMongoDBCache[key]
	db := client.Database(database)
	collection := db.Collection(table)

	return collection
}

func (t *Tools) HandleMongoDBClient() {
	config := config.GetMongodbConfig()
	result := NewMongoDB(config)

	dbMongoDBCache = result.Client

	t.Stdout("MongoDB is Connected")
}

// email
func (t *Tools) GetMailClient(key string) gomail.SendCloser {
	return emailClient[key]
}

func (t *Tools) HandleMailClient() {
	config := config.GetMailConfig()
	result := NewEmail(config)

	emailClient = result.Client

	t.Stdout("Mail Dialer is Connected")
}
