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
	Cache      Storage
	DB         Storage
	Kafka      Storage
	Log        *zap.Logger
	Repository repository.Repository
	Tasks      task.Tasks
}

type Storage interface {
	Close()
}

func NewApp() *App {
	cfg := config.Load()

	cache := cache.NewClient(&cfg.Redis)
	db := mysql.NewClient(&cfg.MySQL)
	kafka := kafka.NewClient(&cfg.Kafka)
	fetch := fetch.NewFetch()
	record := log.Load(cfg.App.Name, cfg.App.Mode)

	app := new(App)

	app.CFG = cfg
	app.Log = record
	app.cacheClient(db, cache, kafka)

	app.Repository = repository.NewRepository(cfg.Host, cache, db, fetch, record)
	app.Tasks = task.NewTasks(record, app.Repository)

	return app
}

func (a *App) cacheClient(db, cache, kafka Storage) {
	a.DB = db
	a.Cache = cache
	a.Kafka = kafka
}

func (a *App) Run(ctx context.Context) {
	a.Stdout("Application initialization started in " + a.CFG.App.Mode + " environment")
	a.Stdout(fmt.Sprintf("Application started successfully and running on %v", os.Getpid()))

	a.Tasks.TaskRuning(ctx)
	a.Tasks.TaskWaiting()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Stdout("Service shutdown initiated")

	if stopper, ok := a.Tasks.(interface{ TaskStopped(context.Context) }); ok {
		stopper.TaskStopped(ctx)
	}

	var errs []error

	if a.Cache != nil {
		a.Stdout("Redis connection closed")
		if err := a.closeStorage(a.Cache, "Redis"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.DB != nil {
		a.Stdout("MySQL connection closed")
		if err := a.closeStorage(a.DB, "MySQL"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.Kafka != nil {
		a.Stdout("Kafka connection closed")
		if err := a.closeStorage(a.Kafka, "Kafka"); err != nil {
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
