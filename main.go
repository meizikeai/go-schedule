package main

import (
	"os"

	"go-schedule/crontab"
	"go-schedule/libs/log"
	"go-schedule/libs/tool"
)

var (
	// rules = tool.NewRegexp()
	// share = tool.NewShare()
	tasks = crontab.NewTasks()
	tools = tool.NewTools()
	// units = tool.NewUnits()
	// chaos  = secret.NewSecret()
	daily = log.NewCreateLog()
	// fetch  = models.NewModelsFetch()
	// jwt    = token.NewJsonWebToken()
	// lion   = fetch.NewFetch()
	// logger = log.NewLogger()
	// logic  = controllers.NewLogic()
	// kafkaProducer = tool.NewKafkaProducer()
	// kafkaConsumer = tool.NewKafkaConsumer()
)

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
