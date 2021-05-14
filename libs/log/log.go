package log

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GIN_ENV")

	infoPath := filepath.Join("/data/logs/", app, "/info.log")
	errorPath := filepath.Join("/data/logs/", app, "/error.log")

	if mode == "debug" {
		infoPath = pwd + "/logs/info.log"
		errorPath = pwd + "/logs/error.log"
	}

	// logrus.SetReportCaller(true)

	infoer, _ := rotatelogs.New(
		infoPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(infoPath),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	errorer, _ := rotatelogs.New(
		errorPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(errorPath),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writerMap := lfshook.WriterMap{
		logrus.DebugLevel: infoer,
		logrus.InfoLevel:  infoer,
		logrus.WarnLevel:  infoer,
		logrus.ErrorLevel: errorer,
		logrus.FatalLevel: errorer,
		logrus.PanicLevel: errorer,
	}

	logrus.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
}
