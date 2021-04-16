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
		// 简单输出
		test.OneJob()

		// 随机插入
		// test.TwoJob()

		// 处理数据
		lover.HandleLoverGift()
	})

	defer a.Stop()

	b := tool.HandleCron("*/1 * * * *", func() {
		// 并发等待组
		test.ThreeJob()

		// FourJob
		test.FourJob()

		// FiveJob
		test.FiveJob()

		// SixJob
		test.SixJob()
	})

	defer b.Stop()

	select {}
}
