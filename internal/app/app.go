// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-schedule/internal/app/repository"
	"go-schedule/internal/app/task"
	"go-schedule/internal/config"
	"go-schedule/internal/pkg/database/cache"
	"go-schedule/internal/pkg/database/kafka"
	"go-schedule/internal/pkg/database/mysql"
	"go-schedule/internal/pkg/fetch"
	"go-schedule/internal/pkg/log"

	"go.uber.org/zap"
)

type App struct {
	CFG        *config.Config
	cache      Storage
	db         Storage
	kafka      Storage
	Log        *zap.Logger
	Repository repository.Repository
	Task       task.Task
}

type Storage interface {
	Close()
}

func NewApp() *App {
	cfg := config.Load()

	cache := cache.NewClient(&cfg.Redis)
	db := mysql.NewClient(&cfg.MySQL)
	kafka := kafka.NewClient(&cfg.Kafka)
	fetch := fetch.NewClient()
	record := log.Load(cfg.App.Name, cfg.App.Mode)

	app := new(App)

	app.CFG = cfg
	app.Log = record
	app.cacheClient(db, cache, kafka)

	app.Repository = repository.New(cache, db, fetch, record, cfg.LB)
	app.Task = task.New(record, app.Repository)

	return app
}

func (a *App) cacheClient(db, cache, kafka Storage) {
	a.db = db
	a.cache = cache
	a.kafka = kafka
}

func (a *App) Run() {
	a.Stdout("Application initialization started in " + a.CFG.App.Mode + " environment")
	a.Stdout(fmt.Sprintf("Application started successfully and running on %v", os.Getpid()))

	a.Task.TaskRuning()
	a.Task.TaskWaiting()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Stdout("Service shutdown initiated")

	if stopper, ok := a.Task.(interface{ TaskStopped(context.Context) }); ok {
		stopper.TaskStopped(ctx)
	}

	var errs []error

	if a.cache != nil {
		a.Stdout("Redis connection closed")
		if err := a.closeStorage(a.cache, "Redis"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.db != nil {
		a.Stdout("MySQL connection closed")
		if err := a.closeStorage(a.db, "MySQL"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.kafka != nil {
		a.Stdout("Kafka connection closed")
		if err := a.closeStorage(a.kafka, "Kafka"); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Service connection closed with error -> %v", errs)
	}

	a.Stdout("Service exited")

	return nil
}

func (a *App) Stdout(format string, v ...any) {
	log := fmt.Sprintf("%s %s %s \n", time.Now().Format(time.DateTime), fmt.Sprintf("[%s]", a.CFG.App.Name), format)

	if _, err := fmt.Fprintf(os.Stdout, log, v...); err != nil {
		fmt.Println(log)
	}
}

func (a *App) closeStorage(s Storage, name string) error {
	defer func() {
		if r := recover(); r != nil {
			a.Stdout("Panic when closing -> %s: %v", name, r)
		}
	}()
	s.Close()
	return nil
}
