package main

import (
	"go-schedule/libs/log"
	"go-schedule/libs/tool"
	"go-schedule/task/lover"
	"go-schedule/task/test"
)

func init() {
	tool.HandleZookeeperConfig()

	tool.HandleLocalMysqlConfig()
	tool.HandleLocalRedisConfig()

	tool.HandleMySQLClient()
	tool.HandleRedisClient()

	log.HandleLogger("go-schedule")
}

func main() {
	a := tool.HandleCron("*/1 * * * *", func() {
		test.OneJob()
		// test.TwoJob()
		lover.HandleLoverGift()
	})

	defer a.Stop()

	b := tool.HandleCron("*/1 * * * *", func() {
		test.ThreeJob()
		test.FourJob()
		test.FiveJob()
		test.SixJob()
	})

	defer b.Stop()

	select {}
}
