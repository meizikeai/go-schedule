package crontab

import (
	"context"
	"fmt"
	"slices"
	"time"

	"go-schedule/config"
	"go-schedule/libs/types"
	"go-schedule/repository"

	binlog "github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

var (
	means        = repository.NewBinlogMySQL()
	binlogConfig = config.GetBinlogConfig()

	// diary             = tool.NewCustomLogger("crontab", false)
	binlogEnableTable = config.GetBinlogEnableTable("default")
)

func (t *Tasks) HandleBinlogStream() {
	config := types.BinlogPosition{
		Name: binlogConfig["default"].Name,
		Pos:  binlogConfig["default"].Pos,
	}
	cache := redis.GetBinlogPosition("default")

	if cache.Name != "" {
		config.Name = cache.Name
	}

	if cache.Pos > 0 {
		config.Pos = cache.Pos
	}

	syncer := tools.GetReplicationMySQLClient("default.master")
	streamer, _ := syncer.StartSync(binlog.Position{
		Name: config.Name,
		Pos:  config.Pos,
	})

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		event, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			continue
		}

		nextPos := syncer.GetNextPosition()
		config.Pos = nextPos.Pos

		switch ev := event.Event.(type) {
		case *replication.RowsEvent:
			// fmt.Printf("[Rows]: %s", ev.Table.Table)
			handleBinlogEvent(event)
		case *replication.RotateEvent:
			// fmt.Printf("[Rotate]: %s, %d\n", ev.NextLogName, ev.Position)
			config.Name = string(ev.NextLogName)
			config.Pos = uint32(ev.Position)
		case *replication.QueryEvent:
			// fmt.Printf("[Query]: %s", ev.Query)
		case *replication.TableMapEvent:
			// fmt.Printf("[TableMap]: %s\n", ev.Table)
		case *replication.XIDEvent:
			// fmt.Printf("[XID]: %d\n", ev.XID)
		case *replication.GTIDEvent:
			// fmt.Printf("[GTID]: %v\n", ev.GNO)
		default:
			// fmt.Printf("[Unknown]: %T\n", event.Event)
		}

		redis.SaveBinlogPosition("zookeeper", string(units.MarshalJson(config)))
	}
}

func handleBinlogEvent(event *replication.BinlogEvent) {
	rowsEvent := event.Event.(*replication.RowsEvent)
	schemaName := string(rowsEvent.Table.Schema)
	tableName := string(rowsEvent.Table.Table)

	schemaTable := fmt.Sprintf("%s.%s", schemaName, tableName)

	if slices.Contains(binlogEnableTable, schemaTable) != true {
		// fmt.Printf("Unauthorized database for %s %s\n", schemaName, tableName)
		return
	}

	tableSchema, _ := means.GetTableSchemaFromCache("default", schemaName, tableName)

	if tableSchema == nil {
		// fmt.Printf("Table schema not found for %s %s\n", schemaName, tableName)
		return
	}

	for _, row := range rowsEvent.Rows {
		data := means.GetBinlogFieldAndValeue(tableSchema, row)
		fmt.Println(data)

		switch event.Header.EventType {
		case replication.WRITE_ROWS_EVENTv2:
			GenerateInsertSQL("default.insert", data)
		case replication.UPDATE_ROWS_EVENTv2:
			GenerateUpdateSQL("default.update", data)
		case replication.DELETE_ROWS_EVENTv2:
			GenerateDeleteSQL("default.delete", data)
		}
	}
}

func GenerateInsertSQL(where string, data map[string]string) {}

func GenerateUpdateSQL(where string, data map[string]string) {}

func GenerateDeleteSQL(where string, data map[string]string) {}
