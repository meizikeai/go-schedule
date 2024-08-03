package tool

import (
	"database/sql"
	"time"

	"go-schedule/config"
	"go-schedule/libs/types"

	"github.com/go-sql-driver/mysql"
)

var connMySQL = types.ConnMySQLMax{
	MaxOpenConns:    200,
	MaxIdleConns:    100,
	ConnmaxLifetime: 10,
}
var fullDbMySQL map[string][]*sql.DB

func (t *Tools) GetMySQLClient(key string) *sql.DB {
	result := fullDbMySQL[key]
	count := t.GetRandmod(len(result))

	return result[count]
}

func (t *Tools) HandleMySQLClient() {
	client := make(map[string][]*sql.DB)

	// local := getMySQLConfig()
	local := config.GetMySQLConfig()

	for k, v := range local {
		m := k + ".master"
		s := k + ".slave"

		for _, addr := range v.Master {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			client[m] = append(client[m], clients)
		}

		for _, addr := range v.Slave {
			clients := handleMySQLClient(addr, v.Username, v.Password, v.Database)
			client[s] = append(client[s], clients)
		}
	}

	fullDbMySQL = client

	t.Stdout("MySQL is Connected")
}

// Timeout, read timeout, write timeout defaults to 1s
func createMySQLClient(config types.OutConfMySQL) *sql.DB {
	dsn := createDSN(config.Addr, config.Username, config.Password, config.Database)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(connMySQL.MaxOpenConns)
	db.SetMaxIdleConns(connMySQL.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMySQL.ConnmaxLifetime) * time.Second)

	err = db.Ping()

	if err != nil {
		panic(err)
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

// Please adjust the connection, read, and write timeouts. The default is 1s.
func createDSN(addr, user, passwd, dbname string) string {
	config := mysql.Config{
		User:             user,                           // Username
		Passwd:           passwd,                         // Password (requires User)
		Net:              "tcp",                          // Network type - default: "tcp"
		Addr:             addr,                           // Network address (requires Net)
		DBName:           dbname,                         // Database name
		MaxAllowedPacket: 4194304,                        // Max packet size allowed  - default: 4194304
		Timeout:          time.Duration(5) * time.Second, // Dial timeout
		ReadTimeout:      time.Duration(5) * time.Second, // I/O read timeout
		WriteTimeout:     time.Duration(5) * time.Second, // I/O write timeout

		AllowNativePasswords: true, // Allows the native password authentication method - default: true
		CheckConnLiveness:    true, // Check connections for liveness before using them - default: true
		InterpolateParams:    true, // Interpolate placeholders into query string  - default: false
	}

	return config.FormatDSN()
}

func (t *Tools) CloseMySQL() {
	for _, val := range fullDbMySQL {
		for _, v := range val {
			v.Close()
		}
	}

	t.Stdout("MySQL is Close")
}
