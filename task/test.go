package task

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func HandleRun() {
	pid := os.Getpid()
	log.Info(fmt.Sprintf("process id %v, go-schedule is working...", pid))
	// log.Error("this is a test")
}

func HandleCheck() func() {
	return func() {
		log.Info("cron is running...")
	}
}
