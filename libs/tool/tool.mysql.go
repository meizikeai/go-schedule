package tool

import (
	"database/sql"
	"time"

	"go-schedule/libs/types"

	"github.com/go-sql-driver/mysql"
)

type configMySQL struct {
	MaxOpenConns    int   `json:"max_open_conns"`
	MaxIdleConns    int   `json:"max_idle_conns"`
	ConnmaxLifetime int64 `json:"conn_max_life_time"`
}

type MySQL struct {
	Client map[string][]*sql.DB
}

var option = configMySQL{
	MaxOpenConns:    200,
	MaxIdleConns:    100,
	ConnmaxLifetime: 10,
}

func NewMySQLClient(data map[string]types.ConfMySQL) *MySQL {
	client := make(map[string][]*sql.DB, 0)

	for k, v := range data {
		m := k + ".master"
		s := k + ".slave"

		for _, addr := range v.Master {
			clients := createMySQLClient(addr, v.Username, v.Password, v.Database)
			client[m] = append(client[m], clients)
		}

		for _, addr := range v.Slave {
			clients := createMySQLClient(addr, v.Username, v.Password, v.Database)
			client[s] = append(client[s], clients)
		}
	}

	return &MySQL{
		Client: client,
	}
}

// Timeout, read timeout, write timeout defaults to 1s
func createMySQLClient(addr, username, password, database string) *sql.DB {
	dsn := createDSN(addr, username, password, database)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(option.MaxOpenConns)
	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(option.ConnmaxLifetime) * time.Second)

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
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

func (m *MySQL) Close() {
	for _, val := range m.Client {
		for _, v := range val {
			v.Close()
		}
	}
}
