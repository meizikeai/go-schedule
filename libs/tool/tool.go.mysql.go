package tool

import (
	"go-schedule/libs/types"

	"github.com/go-mysql-org/go-mysql/replication"
)

type ReplicationMySQL struct {
	Client map[string]*replication.BinlogSyncer
}

func NewReplicationMySQL(data map[string]types.ConfBinlog) *ReplicationMySQL {
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
		ServerID: data.ServerID,
		Flavor:   data.Flavor,
		Host:     data.Host,
		Port:     data.Port,
		User:     data.User,
		Password: data.Password,
	}
	syncer := replication.NewBinlogSyncer(config)
	return syncer
}
