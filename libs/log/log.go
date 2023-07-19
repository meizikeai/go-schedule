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
		MaxSize:    100,   // 日志文件在轮转之前的最大大小，默认 100 MB
		MaxBackups: 10,    // 保留旧日志文件的最大数量
		MaxAge:     15,    // 保留旧日志文件的最大天数
		Compress:   false, // 是否使用 gzip 对日志文件进行压缩归档
		LocalTime:  true,  // 是否使用本地时间，默认 UTC 时间
	}

	return template
}

func createHook(outFile, errFile string) *Hook {
	outlog := getLogger(outFile)
	errlog := getLogger(errFile)

	hook := Hook{
		defaultLogger: outlog,
		minLevel:      logrus.InfoLevel,
		formatter:     &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		loggerByLevel: map[logrus.Level]*lumberjack.Logger{
			logrus.ErrorLevel: errlog,
		},
	}

	return &hook
}

func HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	outFile := filepath.Join("/data/logs/", app, "/out.log")
	errFile := filepath.Join("/data/logs/", app, "/error.log")

	if mode == "debug" {
		outFile = pwd + "/logs/out.log"
		errFile = pwd + "/logs/error.log"
	}

	hook := createHook(outFile, errFile)

	logrus.AddHook(hook)
}
