package tool

import (
	"go-schedule/libs/types"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
)

type ReplicationMySQL struct {
	Client map[string]*replication.BinlogSyncer
}

func NewReplicationMySQLClient(data map[string]types.ConfBinlog) *ReplicationMySQL {
	client := make(map[string]*replication.BinlogSyncer, 0)

	for k, v := range data {
		client[k+".master"] = createReplicationMySQLClient(v)
	}

	return &ReplicationMySQL{
		Client: client,
	}
}

func createReplicationMySQLClient(data types.ConfBinlog) *replication.BinlogSyncer {
	config := replication.BinlogSyncerConfig{
		ServerID: data.ServerId,
		Flavor:   data.Flavor,
		Host:     data.Host,
		Port:     data.Port,
		User:     data.Username,
		Password: data.Password,
	}
	syncer := replication.NewBinlogSyncer(config)
	return syncer
}

type CanalMySQL struct {
	Client map[string]*canal.Canal
}

func NewCanalMySQLClient(data map[string]types.ConfMySQLCanal) *CanalMySQL {
	client := make(map[string]*canal.Canal, 0)

	for k, v := range data {
		client[k+".master"] = createCanalMySQLClient(v)
	}

	return &CanalMySQL{
		Client: client,
	}
}

func createCanalMySQLClient(data types.ConfMySQLCanal) *canal.Canal {
	config := canal.NewDefaultConfig()

	config.Addr = data.Addr
	config.User = data.Username
	config.Password = data.Password

	config.ServerID = data.ServerId
	config.Flavor = data.Flavor

	config.Dump.TableDB = "test"
	config.Dump.Tables = []string{"canal_test"}

	c, err := canal.NewCanal(config)

	if err != nil {
		panic(err)
	}

	return c
}
