// internal/pkg/database/mysql/client.go
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"go-schedule/internal/config"

	"github.com/go-sql-driver/mysql"
)

type Clients struct {
	clients map[string][]*sql.DB
	mu      sync.RWMutex
}

func NewClient(cfg *map[string][]config.MySQLInstance) *Clients {
	c := &Clients{
		clients: make(map[string][]*sql.DB),
	}

	for key, value := range *cfg {
		for _, v := range value {
			for _, dsn := range v.Master {
				db := createClient(dsn, &v)
				if db != nil {
					c.clients[key] = append(c.clients[key], db)
				}
			}
			for _, dsn := range v.Slave {
				db := createClient(dsn, &v)
				if db != nil {
					slaveKey := fmt.Sprintf("%s.slave", key)
					c.clients[slaveKey] = append(c.clients[slaveKey], db)
				}
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
	cfg.InterpolateParams = true

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
	c.mu.RLock()
	defer c.mu.RUnlock()

	if clients, ok := c.clients[key]; ok && len(clients) > 0 {
		return clients[rand.Intn(len(clients))]
	}

	if strings.HasSuffix(key, ".slave") {
		masterKey := key[:len(key)-6]
		if clients, ok := c.clients[masterKey]; ok && len(clients) > 0 {
			return clients[rand.Intn(len(clients))]
		}
	}

	return nil
}

func (c *Clients) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, value := range c.clients {
		for _, v := range value {
			v.Close()
		}
	}
}
