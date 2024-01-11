package task

import (
	"fmt"
	"os"

	"go-schedule/config"
	"go-schedule/libs/tool"
)

func HandleRun() {
	pid := os.Getpid()
	tool.Stdout("The current environment is " + config.GetMode())
	tool.Stdout(fmt.Sprintf("The process id of the service is %v", pid))
}

func HandleCheck() func() {
	return func() {
		tool.Stdout("Scheduled task is running...")
	}
}
