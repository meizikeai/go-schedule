package main

import (
	"os"

	"go-schedule/libs/log"
	"go-schedule/libs/tool"
	"go-schedule/task"
)

func init() {
	// tool.HandleMySQLClient()
	// tool.HandleRedisClient()
	// tool.HandleMongodbClient()

	// tool.HandleElasticSearchClient()
	// tool.HandleKafkaProducerClient()
	// tool.HandleKafkaConsumerClient()

	log.HandleLogger("go-schedule")
}

func main() {
	tool.SignalHandler(func() {
		tool.CloseMySQL()
		tool.CloseRedis()
		// tool.CloseMongoDB()
		// tool.CloseKafka()
		// tool.CloseElasticSearch()

		tool.Stdout("Server Shutdown")

		os.Exit(0)
	})

	// running
	task.HandleRun()

	// checking
	check := tool.HandleCron("0 10 */1 * *", task.HandleCheck())
	defer check.Stop()

	// kafka producer
	// tool.SendKafkaProducerMessage("broker", "topic", "sync", "test")

	// kafka consumer
	// tool.HandlerKafkaConsumerMessage("broker", "topic")

	select {}
}
