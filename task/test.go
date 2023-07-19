package task

import (
	log "github.com/sirupsen/logrus"
)

func HandleRun() {
	log.Info("go-schedule is working...")
	// log.Error("this is a test")
}

func HandleCheck() func() {
	return func() {
		log.Info("cron is running...")
	}
}
