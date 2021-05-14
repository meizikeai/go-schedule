package tool

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go-schedule/libs/types"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var fullDbMySQL map[string]*sql.DB
var mysqlConfig types.FullConfMySQL

func HandleLocalMysqlConfig() {
	var config types.FullConfMySQL

	pwd, _ := os.Getwd()
	mode := GetGinMODE()

	address := strings.Join([]string{pwd, "/conf/", mode, ".mysql.json"}, "")

	res, err := ioutil.ReadFile(address)

	if err != nil {
		log.Fatal(err)
	}

	// 解析json
	err = json.Unmarshal(res, &config)

	if err != nil {
		log.Fatal(err)
	}

	mysqlConfig = config
}

func GetMySQLClient(key string) *sql.DB {
	return fullDbMySQL[key]
}

func HandleMySQLClient() {
	config := make(map[string]*sql.DB)

	one := getZookeeperMysqlConfig()
	two := getLocalMysqlConfig()

	for key, val := range one {
		m, s := getMySQLConfig(val)

		config[key+".master"] = createMySQLClient(m)
		config[key+".slave"] = createMySQLClient(s)
	}

	for key, val := range two {
		m, s := getMySQLConfig(val)

		config[key+".master"] = createMySQLClient(m)
		config[key+".slave"] = createMySQLClient(s)
	}

	fullDbMySQL = config
}

func getLocalMysqlConfig() types.FullConfMySQL {
	return mysqlConfig
}

func createMySQLClient(config types.OutConfMySQL) *sql.DB {
	// db *sql.DB
	// 可长期存在，不建议频繁开启/关闭
	// defer db.Close()

	path := strings.Join([]string{config.Username, ":", config.Password, "@tcp(", config.Addr, ")/", config.Database, "?charset=utf8"}, "")

	// create client
	db, err := sql.Open("mysql", path)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	// test client
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getMySQLConfig(config types.ConfMySQL) (types.OutConfMySQL, types.OutConfMySQL) {
	master := config.Master
	slave := config.Slave

	i := GetRandmod(len(master))
	j := GetRandmod(len(slave))

	m := types.OutConfMySQL{
		Addr:     master[i],
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	}

	s := types.OutConfMySQL{
		Addr:     slave[j],
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	}

	return m, s
}
