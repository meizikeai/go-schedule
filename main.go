package main

import (
	"go-schedule/libs/log"
	"go-schedule/libs/tool"
	"go-schedule/task/test"
)

func init() {
	tool.HandleZookeeperConfig()

	// not recommended for use
	// tool.HandleLocalMysqlConfig()
	// tool.HandleLocalRedisConfig()

	tool.HandleMySQLClient()
	tool.HandleRedisClient()
	// tool.HandleKafkaProducerClient()
	// tool.HandleKafkaConsumerClient()

	log.HandleLogger("go-schedule")
}

func main() {
	a := tool.HandleCron("*/1 * * * *", func() {
		test.OneJob()
		// test.TwoJob()
		// test.HandleLoverGift()
	})

	defer a.Stop()

	b := tool.HandleCron("*/1 * * * *", func() {
		test.ThreeJob()
		test.FourJob()
		// test.FiveJob()
		// test.SixJob()
	})

	defer b.Stop()

	// kafka producer
	// tool.SendKafkaProducerMessage("broker", "topic", "sync", "test")

	// kafka consumer
	// tool.HandlerKafkaConsumerMessage("broker", "topic")

	select {}
}
