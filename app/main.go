package main

import (
	"os"

	"go-schedule/crontab"
	"go-schedule/libs/tool"
)

var (
	// chaos  = tool.NewSecret()
	daily = tool.NewCreateLog()
	// jwt    = tool.NewJsonWebToken()
	lion   = tool.NewFetch()
	logger = tool.NewLogger()
	// rules  = tool.NewRegexp()
	// share  = tool.NewShare()
	tools = tool.NewTools()
	units = tool.NewUnits()
	// logic = controllers.NewLogic()
	tasks = crontab.NewTasks()
	// fetch  = models.NewModelsFetch()
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
	// tools.HandleKafkaConsumerGroupClient()

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

	select {}
}
