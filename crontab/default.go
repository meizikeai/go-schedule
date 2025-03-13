package crontab

import (
	"fmt"
	"os"

	"go-schedule/config"
	"go-schedule/libs/tool"
	"go-schedule/models"
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

var (
	tools = tool.NewTools()
	units = tool.NewUnits()
	mysql = models.NewModelsMySQL()
	redis = models.NewModelsRedis()
)

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
