package task

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func HandleRun() {
	pid := os.Getpid()
	log.Info(fmt.Sprintf("The process id of the service is %v", pid))
}

func HandleCheck() func() {
	return func() {
		log.Info("Scheduled task is running...")
	}
}
