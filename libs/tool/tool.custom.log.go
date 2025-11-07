package tool

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func NewCustomLogger(name string) (*slog.Logger, error) {
	basePath := getLogPath(name)

	file, err := os.OpenFile(basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.StringValue(time.Now().Format("2006-01-02 15:04:05"))
			case slog.LevelKey:
				return slog.Attr{}
			case slog.SourceKey:
				if src, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(filepath.Base(src.File) + ":" + strconv.Itoa(src.Line))
				}
			}
			return a
		},
	})

	return slog.New(handler), nil
}

func getLogPath(name string) string {
	appName := os.Getenv("GO_APP")
	if appName == "" {
		appName = "custom-log"
	}

	var logDir string
	if os.Getenv("GO_ENV") == "debug" {
		pwd, _ := os.Getwd()
		logDir = filepath.Join(pwd, "../logs")
	} else {
		logDir = filepath.Join("/data/logs", appName)
	}

	_ = os.MkdirAll(logDir, 0755)

	timestamp := "." + time.Now().Format("2006.01.02.150405")

	return filepath.Join(logDir, fmt.Sprintf("%s.log%s", name, timestamp))
}
