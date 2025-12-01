package types

type ConfBinlog struct {
	ServerID uint32 `json:"serverId"`
	Flavor   string `json:"flavor"`
	Name     string `json:"name"`
	Pos      uint32 `json:"pos"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ConfMySQLCanal struct {
	ServerId uint32 `json:"serverId"`
	Flavor   string `json:"flavor"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TableColumnRows struct {
	Name      string
	Type      string
	Collation string
	Extra     string
}

type BinlogPosition struct {
	Name string `json:"name"`
	Pos  uint32 `json:"pos"`
}
