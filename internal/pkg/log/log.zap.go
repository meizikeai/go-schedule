// internal/pkg/log/log.zap.go
package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Sugar *zap.SugaredLogger
)

func Load(app, env string) *zap.Logger {
	atomicLevel := zap.NewAtomicLevelAt(getLevel(env))

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("/data/logs/%s/app.log", app),
		MaxSize:    500, // Maximum log file split size, default 500 MB
		MaxBackups: 50,  // Maximum number of old log files to keep
		MaxAge:     30,  // Maximum number of days to keep old log files
		Compress:   true,
	})

	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("/data/logs/%s/error.log", app),
		MaxSize:    100,
		MaxBackups: 50,
		MaxAge:     60,
		Compress:   true,
	})

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	})

	core := zapcore.NewTee(
		// zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer), atomicLevel),
		zapcore.NewCore(encoder, writer, atomicLevel),
		zapcore.NewCore(encoder, errorWriter, zapcore.ErrorLevel),
	)

	L := zap.New(core, zap.AddCaller())

	Sugar = L.Sugar()
	zap.ReplaceGlobals(L)

	return L
}

func getLevel(env string) zapcore.Level {
	switch env {
	case "local", "debug":
		return zap.DebugLevel
	case "test":
		return zap.InfoLevel
	default:
		return zap.InfoLevel
	}
}
