package types

type ConfMySQL struct {
	Master   []string `json:"master"`
	Slave    []string `json:"slave"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Database string   `json:"database"`
}

type ConfRedis struct {
	Master   []string `json:"master"`
	Password string   `json:"password"`
	Db       int      `json:"db"`
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

type ConfMail struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfigKafkaConsumerCroup struct {
	Assignor string `json:"assignor"`
	GroupId  string `json:"group_id"`
	Oldest   bool   `json:"oldest"`
	Version  string `json:"version"`
}
