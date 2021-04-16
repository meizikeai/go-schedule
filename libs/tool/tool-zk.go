package tool

import (
	"encoding/json"
	"go-schedule/conf"
	"go-schedule/libs/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
)

type confZookeeper map[string]map[string]string

var zookeeperMySQL types.FullConfMySQL
var zookeeperRedis types.FullConfRedis
var zookeeperServer types.MapStringString
var zookeeperString types.MapStringString

func HandleZookeeperConfig() {
	zkList := conf.Release

	mode := os.Getenv("GIN_MODE")

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
		} else {
			config := make(types.MapStringString)

			for key, val := range v {
				back := getServerZookeeperChildren(pool, key, val)
				config[key] = back
			}

			zookeeperServer = config
		}
	}
	log.Info(zookeeperMySQL)
	log.Info(zookeeperRedis)
	log.Info(zookeeperServer)
}

func GetZookeeperServerConfig() types.MapStringString {
	return zookeeperServer
}

func GetZookeeperRedisConfig() types.MapStringString {
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
	// address := strings.Join([]string{pwd, "/conf/zk.json"}, "")
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

func getServerZookeeperChildren(pool *zk.Conn, key string, path string) string {
	var result string

	back, _, err := pool.Children(path)

	if err != nil {
		log.Fatal(err)
	}

	i := GetRandmod(len(back))
	result = back[i]

	return result
}
