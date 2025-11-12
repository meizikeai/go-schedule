package crontab

import (
	"fmt"
	"os"

	"go-schedule/config"
	"go-schedule/libs/tool"
	"go-schedule/repository"
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

var (
	tools = tool.NewTools()
	units = tool.NewUnits()
	// mysql = repository.NewModelsMySQL()
	redis = repository.NewModelsRedis()
)

func (t *Tasks) HandleRun() {
	pid := os.Getpid()
	tools.Stdout("Starting application in the " + config.GetMode() + " environment")
	tools.Stdout(fmt.Sprintf("Service process ID %v", pid))
}

func (t *Tasks) HandleCheck() func() {
	return func() {
		tools.Stdout("Scheduled tasks are running...")
	}
}
