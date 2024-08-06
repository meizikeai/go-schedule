package main

import (
	"os"

	"go-schedule/crontab"
	"go-schedule/libs/log"
	"go-schedule/libs/tool"
)

var tools = tool.NewTools()
var tasks = crontab.NewTasks()
var daily = log.NewCreateLog()

// var kafkaProducer = tool.NewKafkaProducer()
// var kafkaConsumer = tool.NewKafkaConsumer()

func init() {
	// tools.HandleZookeeperClient()
	// tools.HandleMySQLClient()
	// tools.HandleRedisClient()
	// tools.HandleMongodbClient()
	// tools.HandleMailClient()

	// tools.HandleElasticSearchClient()
	// tools.HandleKafkaProducerClient()
	// tools.HandleKafkaConsumerClient()

	daily.HandleLogger("go-schedule")
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
	tasks.HandleRun()

	// checking
	check := tools.HandleCron("0 10 */1 * *", tasks.HandleCheck())
	defer check.Stop()

	// kafka producer
	// kafkaProducer.SendKafkaProducerMessage("broker", "topic", "sync", "test")

	// kafka consumer
	// kafkaConsumer.HandlerKafkaConsumerMessage("broker", "topic")

	select {}
}
