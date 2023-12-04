package log

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Hook struct {
	defaultLogger *lumberjack.Logger
	formatter     logrus.Formatter
	minLevel      logrus.Level
	loggerByLevel map[logrus.Level]*lumberjack.Logger
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)

	if err != nil {
		return err
	}

	if logger, ok := hook.loggerByLevel[entry.Level]; ok {
		_, err = logger.Write([]byte(msg))
	} else {
		_, err = hook.defaultLogger.Write([]byte(msg))
	}

	return err
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.minLevel+1]
}

func getLogger(file string) *lumberjack.Logger {
	template := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100,   // Maximum log file split size, default 100 MB
		MaxBackups: 10,    // Maximum number of old log files to keep
		MaxAge:     15,    // Maximum number of days to keep old log files
		Compress:   false, // Whether to use gzip to compress and archive log files
		LocalTime:  true,  // Whether to use local time, default UTC time
	}

	return template
}

func createHook(traFile, outFile, errFile string) *Hook {
	tralog := getLogger(traFile)
	outlog := getLogger(outFile)
	errlog := getLogger(errFile)

	hook := Hook{
		defaultLogger: outlog,
		minLevel:      logrus.TraceLevel,
		formatter:     &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		loggerByLevel: map[logrus.Level]*lumberjack.Logger{
			logrus.TraceLevel: tralog,
			logrus.ErrorLevel: errlog,
		},
	}

	return &hook
}

func HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	traFile := filepath.Join("/data/logs/", app, "/trace.log")
	outFile := filepath.Join("/data/logs/", app, "/out.log")
	errFile := filepath.Join("/data/logs/", app, "/error.log")

	if mode == "debug" {
		traFile = pwd + "/logs/trace.log"
		outFile = pwd + "/logs/out.log"
		errFile = pwd + "/logs/error.log"
	}

	hook := createHook(traFile, outFile, errFile)

	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(hook)
}
