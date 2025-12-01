// internal/app/task/default.go
package task

import (
	"context"

	"go-schedule/internal/app/repository"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Tasks interface {
	TaskRuning()
	TaskStopped(ctx context.Context)
	TaskWaiting()
}

type tasks struct {
	log     *zap.Logger
	repo    repository.Repository
	cron    *cron.Cron
	started chan struct{}
}

func NewTasks(log *zap.Logger, repo repository.Repository) Tasks {
	return &tasks{
		log:     log,
		repo:    repo,
		cron:    cron.New(),
		started: make(chan struct{}),
	}
}

func (t *tasks) TaskRuning() {
	// every day 10:00
	t.cron.AddFunc("0 10 */1 * *", func() {
		t.log.Info("Cron task always running")
	})

	// todo
	// t.cron.AddFunc("@every 10s", func() { ... })

	t.cron.Start()
	close(t.started)
}

func (t *tasks) TaskWaiting() {
	<-t.started
}

func (t *tasks) TaskStopped(ctx context.Context) {
	cronCtx := t.cron.Stop()

	select {
	case <-cronCtx.Done():
		t.log.Info("All cron tasks stopped cleanly")
	case <-ctx.Done():
		t.log.Warn("Shutdown timeout: cron tasks may still be running", zap.Error(ctx.Err()))
	}
}
