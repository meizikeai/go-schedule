package tool

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go-schedule/libs/types"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var connMySQL = types.ConnMySQLMax{
	MaxOpenConns:    500,
	MaxIdleConns:    250,
	ConnmaxLifetime: 10,
}
var fullDbMySQL map[string][]*sql.DB
var mysqlConfig types.FullConfMySQL

func HandleLocalMysqlConfig() {
	var config types.FullConfMySQL

	pwd, _ := os.Getwd()
	mode := GetMODE()

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
	result := fullDbMySQL[key]
	count := GetRandmod(len(result))

	return result[count]
}

func HandleMySQLClient() {
	config := make(map[string][]*sql.DB)

	zookeeper := getZookeeperMysqlConfig()
	local := getLocalMysqlConfig()

	for k, v := range zookeeper {
		m := k + ".master"
		s := k + ".slave"

		for _, addr := range v.Master {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			config[m] = append(config[m], clients)
		}

		for _, addr := range v.Slave {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			config[s] = append(config[s], clients)
		}
	}

	for k, v := range local {
		m := k + ".master"
		s := k + ".slave"

		for _, addr := range v.Master {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			config[m] = append(config[m], clients)
		}

		for _, addr := range v.Slave {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			config[s] = append(config[s], clients)
		}
	}

	fullDbMySQL = config
}

func getLocalMysqlConfig() types.FullConfMySQL {
	return mysqlConfig
}

// timeout、readTimeout、writeTimeout default 1s
func createMySQLClient(config types.OutConfMySQL) *sql.DB {
	dsn := createDSN(config.Addr, config.Username, config.Password, config.Database)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(connMySQL.MaxOpenConns)
	db.SetMaxIdleConns(connMySQL.MaxIdleConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(connMySQL.ConnmaxLifetime))

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func handleMySQLClient(addr, username, password, database string) *sql.DB {
	option := types.OutConfMySQL{
		Addr:     addr,
		Username: username,
		Password: password,
		Database: database,
	}

	client := createMySQLClient(option)

	return client
}

// 连接、读、写超时请安需调整，默认一秒
func createDSN(addr, user, passwd, dbname string) string {
	config := mysql.Config{
		User:             user,                            // Username
		Passwd:           passwd,                          // Password (requires User)
		Net:              "tcp",                           // Network type - default: "tcp"
		Addr:             addr,                            // Network address (requires Net)
		DBName:           dbname,                          // Database name
		MaxAllowedPacket: 4194304,                         // Max packet size allowed  - default: 4194304
		Timeout:          time.Second * time.Duration(10), // Dial timeout
		ReadTimeout:      time.Second * time.Duration(30), // I/O read timeout
		WriteTimeout:     time.Second * time.Duration(60), // I/O write timeout

		AllowNativePasswords: true, // Allows the native password authentication method - default: true
		CheckConnLiveness:    true, // Check connections for liveness before using them - default: true
		InterpolateParams:    true, // Interpolate placeholders into query string  - default: false
	}

	return config.FormatDSN()
}

func CloseMySQL() {
	for _, val := range fullDbMySQL {
		for _, v := range val {
			v.Close()
		}
	}

	Stdout("MySQL Close")
}
