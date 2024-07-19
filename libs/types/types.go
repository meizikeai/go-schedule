package types

type MapStringInterface map[string]any
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

type ConfRedis struct {
	Master   []string `json:"master"`
	Password string   `json:"password"`
	Db       int      `json:"db"`
}
type OutConfRedis struct {
	Addr               string `json:"addr"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	DB                 int    `json:"db"`
	MaxRetries         int    `json:"max_retries"`
	PoolSize           int    `json:"pool_size"`
	MinIdleConns       int    `json:"min_idle_conns"`
	DialTimeout        int    `json:"dial_timeout"`
	ReadTimeout        int    `json:"read_timeout"`
	WriteTimeout       int    `json:"write_timeout"`
	IdleTimeout        int    `json:"idle_timeout"`
	IdleCheckFrequency int    `json:"idle_check_frequency"`
}

type ConnMySQLMax struct {
	MaxOpenConns    int   `json:"max_open_conns"`
	MaxIdleConns    int   `json:"max_idle_conns"`
	ConnmaxLifetime int64 `json:"conn_max_life_time"`
}

type ConfElasticSearch struct {
	Address  []string `json:"address"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

type ConfEtcd struct {
	Address  []string `json:"address"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

type ConfMongoDB struct {
	Master string `json:"master"`
	Slave  string `json:"slave"`
}

type ConMail struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MailMessage struct {
	From    string   `json:"from"`
	Subject string   `json:"subject"`
	Data    string   `json:"data"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	File    []string `json:"file"`
}
