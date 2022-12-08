package tool

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-schedule/conf"
	"go-schedule/libs/types"

	"github.com/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
)

type confZookeeper map[string]map[string]string

var zookeeperMySQL types.FullConfMySQL
var zookeeperRedis types.FullConfRedis
var zookeeperServer types.MapStringString
var zookeeperString types.MapStringString

func HandleZookeeperConfig() {
	zkList := conf.Release

	mode := GetMODE()

	if mode != "release" {
		zkList = conf.Test
	}

	pool := newZookeeperConnect(zkList)

	defer pool.Close()

	option := readZookeeperConfig()

	for k, v := range option {
		if k == "mysql" {
			config := make(types.FullConfMySQL)

			for key, val := range v {
				back := getMysqlZookeeperChildren(pool, val)
				config[key] = back
			}

			zookeeperMySQL = config
		} else if k == "redis" {
			config := make(types.FullConfRedis)

			for key, val := range v {
				back := getRedisZookeeperChildren(pool, val)
				config[key] = back
			}

			zookeeperRedis = config
		} else if k == "string" {
			config := make(types.MapStringString)

			for key, val := range v {
				back := getZookeeperGet(pool, val)
				config[key] = back
			}

			zookeeperString = config
		} else if k == "kafka" {
			config := make(types.MapStringString)

			for key, val := range v {
				back := getKafkaZookeeperChildren(pool, key, val)
				config[key] = back
			}

			zookeeperString = config
		} else {
			config := make(types.MapStringString)

			for key, val := range v {
				back := getServerZookeeperChildren(pool, key, val)
				config[key] = back
			}

			zookeeperServer = config
		}
	}
}

func GetZookeeperServerConfig() types.MapStringString {
	return zookeeperServer
}

func GetZookeeperStringConfig() types.MapStringString {
	return zookeeperString
}

func getZookeeperMysqlConfig() types.FullConfMySQL {
	return zookeeperMySQL
}

func getZookeeperRedisConfig() types.FullConfRedis {
	return zookeeperRedis
}

func readZookeeperConfig() confZookeeper {
	var config confZookeeper

	pwd, _ := os.Getwd()
	address := filepath.Join(pwd, "/conf/zk.json")
	res, err := ioutil.ReadFile(address)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func newZookeeperConnect(zkList []string) (pool *zk.Conn) {
	pool, _, err := zk.Connect(zkList, 20*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func getZookeeperGet(pool *zk.Conn, path string) string {
	v, _, err := pool.Get(path)

	if err != nil {
		log.Fatal(err)
	}

	return string(v)
}

func getMysqlZookeeperChildren(pool *zk.Conn, path string) types.ConfMySQL {
	var result types.ConfMySQL

	mysql := types.ConfMySQL{
		Master:   nil,
		Slave:    nil,
		Username: "",
		Password: "",
		Database: "",
	}

	back, _, err := pool.Children(path)

	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile(`:`)

	for _, val := range back {
		res := reg.FindString(val)

		if res != ":" {
			switch val {
			case "master":
				res, _, err := pool.Children(path + "/" + val)

				if err != nil {
					log.Fatal(err)
				}

				mysql.Master = res
			case "slave":
				res, _, err := pool.Children(path + "/" + val)

				if err != nil {
					log.Fatal(err)
				}

				mysql.Slave = res
			case "username":
				mysql.Username = getZookeeperGet(pool, path+"/"+val)
			case "password":
				mysql.Password = getZookeeperGet(pool, path+"/"+val)
			case "database":
				mysql.Database = getZookeeperGet(pool, path+"/"+val)
			}

			result = mysql
		}
	}

	return result
}

func getRedisZookeeperChildren(pool *zk.Conn, path string) types.ConfRedis {
	var result types.ConfRedis

	redis := types.ConfRedis{
		Master:   nil,
		Password: "",
		Db:       0,
	}

	back, _, err := pool.Children(path)

	if err != nil {
		log.Fatal(err)
	}

	redis.Master = back
	result = redis

	return result
}

func getServerZookeeperChildren(pool *zk.Conn, key, path string) string {
	var result string

	back, _, err := pool.Children(path)

	if err != nil {
		log.Fatal(err)
	}

	i := GetRandmod(len(back))
	result = back[i]

	return result
}

func getKafkaZookeeperChildren(pool *zk.Conn, key, path string) string {
	var result string

	back, _, err := pool.Get(path)

	if err != nil {
		log.Fatal(err)
	}

	res := string(back)
	list := strings.Split(res, ",")

	for k, v := range list {
		if k+1 == len(list) {
			result += v + ":9092"
		} else {
			result += v + ":9092,"
		}
	}

	return result
}
