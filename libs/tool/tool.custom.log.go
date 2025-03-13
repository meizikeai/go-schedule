package tool

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func NewCustomLogger(name string, fixed bool) *slog.Logger {
	path := fmt.Sprintf("/data/logs/%s/", GetAppName())

	if GetGoEnv() == "debug" {
		pwd, _ := os.Getwd()
		path = filepath.Join(pwd, "../logs/")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil
	}

	trace := ""

	if fixed == true {
		trace = time.Now().Format(".2006.01.02.150405")
	}

	basePath := filepath.Join(path, fmt.Sprintf("%s.log%s", name, trace))
	file, err := os.OpenFile(basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil
	}

	h := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().Format("2006-01-02 15:04:05"))
			}
			if a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(filepath.Base(src.File) + ":" + strconv.Itoa(src.Line))
				}
			}
			return a
		},
	})

	return slog.New(h)
}
