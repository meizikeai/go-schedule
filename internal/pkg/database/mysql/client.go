// internal/pkg/database/mysql/client.go
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-schedule/internal/config"

	"github.com/go-sql-driver/mysql"
)

type Clients struct {
	clients map[string][]*sql.DB
	rand    *rand.Rand
}

func NewClient(cfg *map[string][]config.MySQLInstance) *Clients {
	c := &Clients{
		clients: make(map[string][]*sql.DB),
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for key, value := range *cfg {
		for _, v := range value {
			for _, dsn := range v.Master {
				db := createClient(dsn, &v)
				c.clients[key] = append(c.clients[key], db)
			}
			for _, dsn := range v.Slave {
				db := createClient(dsn, &v)
				slaveKey := fmt.Sprintf("%s.slave", key)
				c.clients[slaveKey] = append(c.clients[slaveKey], db)
			}
		}
	}

	return c
}

func createClient(dsn string, option *config.MySQLInstance) *sql.DB {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	setConnPool(db, option)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	return db
}

func setConnPool(db *sql.DB, cfg *config.MySQLInstance) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
}

func (c *Clients) Client(key string) *sql.DB {
	if clients := c.clients[key]; len(clients) > 0 {
		return clients[c.rand.Intn(len(clients))]
	}

	if strings.HasSuffix(key, ".slave") {
		masterKey := key[:len(key)-6]
		if clients := c.clients[masterKey]; len(clients) > 0 {
			return clients[c.rand.Intn(len(clients))]
		}
	}

	return nil
}

func (c *Clients) Close() {
	for _, value := range c.clients {
		for _, v := range value {
			v.Close()
		}
	}
}
