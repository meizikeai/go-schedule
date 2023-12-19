package task

import (
	"fmt"
	"os"

	"go-schedule/libs/tool"
)

func HandleRun() {
	pid := os.Getpid()
	tool.Stdout(fmt.Sprintf("The process id of the service is %v", pid))
}

func HandleCheck() func() {
	return func() {
		tool.Stdout("Scheduled task is running...")
	}
}
