package types

type MapStringInterface map[string]interface{}
type MapStringString map[string]string

type ConfMySQL struct {
	Master   []string `json:"master"`
	Slave    []string `json:"slave"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Database string   `json:"database"`
}
type OutConfMySQL struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
type FullConfMySQL map[string]ConfMySQL

type ConfRedis struct {
	Master   []string `json:"master"`
	Password string   `json:"password"`
	Db       int      `json:"db"`
}
type OutConfRedis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`

	MaxRetries         int `json:"max_retries"`
	PoolSize           int `json:"pool_size"`
	ReadTimeout        int `json:"read_timeout"`
	WriteTimeout       int `json:"write_timeout"`
	IdleTimeout        int `json:"idle_timeout"`
	IdleCheckFrequency int `json:"idle_check_frequency"`
	MaxConnAge         int `json:"max_conn_age"`
	PoolTimeout        int `json:"pool_timeout"`
}
type FullConfRedis map[string]ConfRedis

type ConnMySQLMax struct {
	MaxLifetime int64 `json:"MaxLifetime"`
	MaxIdleConn int   `json:"MaxIdleConn"`
	MaxOpenConn int   `json:"MaxOpenConn"`
}
