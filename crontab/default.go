package crontab

import (
	"fmt"
	"os"

	"go-schedule/config"
	"go-schedule/libs/tool"
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

var tools = tool.NewTools()

func (t *Tasks) HandleRun() {
	pid := os.Getpid()
	tools.Stdout("The current environment is " + config.GetMode())
	tools.Stdout(fmt.Sprintf("The process id of the service is %v", pid))
}

func (t *Tasks) HandleCheck() func() {
	return func() {
		tools.Stdout("Scheduled task is running...")
	}
}
