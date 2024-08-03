package task

import (
	"fmt"
	"os"

	"go-schedule/config"
	"go-schedule/libs/tool"
)

var tools = tool.NewTools()

func HandleRun() {
	pid := os.Getpid()
	tools.Stdout("The current environment is " + config.GetMode())
	tools.Stdout(fmt.Sprintf("The process id of the service is %v", pid))
}

func HandleCheck() func() {
	return func() {
		tools.Stdout("Scheduled task is running...")
	}
}
