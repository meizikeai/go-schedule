package log

import (
	"io"
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

func createHook(errFile, warFile, infFile, debFile, traFile string) *Hook {
	errlog := getLogger(errFile)
	warlog := getLogger(warFile)
	inflog := getLogger(infFile)
	deblog := getLogger(debFile)
	tralog := getLogger(traFile)

	hook := Hook{
		defaultLogger: tralog,
		minLevel:      logrus.TraceLevel,
		formatter:     &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		loggerByLevel: map[logrus.Level]*lumberjack.Logger{
			logrus.ErrorLevel: errlog,
			logrus.WarnLevel:  warlog,
			logrus.InfoLevel:  inflog,
			logrus.DebugLevel: deblog,
			logrus.TraceLevel: tralog,
		},
	}

	return &hook
}

func HandleLogger(app string) {
	pwd, _ := os.Getwd()
	mode := os.Getenv("GO_ENV")

	errFile := filepath.Join("/data/logs/", app, "/error.log")
	warFile := filepath.Join("/data/logs/", app, "/warn.log")
	infFile := filepath.Join("/data/logs/", app, "/info.log")
	debFile := filepath.Join("/data/logs/", app, "/debug.log")
	traFile := filepath.Join("/data/logs/", app, "/trace.log")

	if mode == "debug" {
		errFile = filepath.Join(pwd, "../logs/error.log")
		warFile = filepath.Join(pwd, "../logs/warn.log")
		infFile = filepath.Join(pwd, "../logs/info.log")
		debFile = filepath.Join(pwd, "../logs/debug.log")
		traFile = filepath.Join(pwd, "../logs/trace.log")
	}

	hook := createHook(errFile, warFile, infFile, debFile, traFile)

	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(hook)
}
