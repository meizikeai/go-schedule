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

var ConnMax = types.ConnMax{
	MaxLifetime: 4,
	MaxIdleConn: 200,
	MaxOpenConn: 200,
}
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
	path := strings.Join([]string{config.Username, ":", config.Password, "@tcp(", config.Addr, ")/", config.Database, "?charset=utf8"}, "")

	db, err := sql.Open("mysql", path)

	db.SetConnMaxLifetime(time.Duration(ConnMax.MaxLifetime) * time.Hour)
	db.SetMaxIdleConns(ConnMax.MaxIdleConn)
	db.SetMaxOpenConns(ConnMax.MaxOpenConn)

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
