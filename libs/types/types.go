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
}
type FullConfRedis map[string]ConfRedis

type ConnMax struct {
	MaxLifetime int64 `json:"MaxLifetime"`
	MaxIdleConn int   `json:"MaxIdleConn"`
	MaxOpenConn int   `json:"MaxOpenConn"`
}
