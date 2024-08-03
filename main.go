package main

import (
	"os"

	"go-schedule/libs/log"
	"go-schedule/libs/tool"
	"go-schedule/task"
)

var tools = tool.NewTools()

func init() {
	// tools.HandleZookeeperClient()
	// tools.HandleMySQLClient()
	// tools.HandleRedisClient()
	// tools.HandleMongodbClient()
	// tools.HandleMailClient()

	// tools.HandleElasticSearchClient()
	// tools.HandleKafkaProducerClient()
	// tools.HandleKafkaConsumerClient()

	log.HandleLogger("go-schedule")
}

func main() {
	tools.SignalHandler(func() {
		// tools.CloseMySQL()
		// tools.CloseRedis()
		// tools.CloseMongoDB()
		// tools.CloseKafka()
		// tools.CloseElasticSearch()
		// tools.CloseMail()

		tools.Stdout("The Service is Shutdown")

		os.Exit(0)
	})

	// running
	task.HandleRun()

	// checking
	check := tools.HandleCron("0 10 */1 * *", task.HandleCheck())
	defer check.Stop()

	// kafka producer
	// tool.SendKafkaProducerMessage("broker", "topic", "sync", "test")

	// kafka consumer
	// tool.HandlerKafkaConsumerMessage("broker", "topic")

	select {}
}
