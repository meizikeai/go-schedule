package tool

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func HandleCron(time string, fn func()) *cron.Cron {
	res := cron.New()

	_, err := res.AddFunc(time, fn)

	if err != nil {
		fmt.Println(err)
	}

	res.Start()

	return res
}
