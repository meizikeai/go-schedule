// cmd/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-schedule/internal/app"
)

func main() {
	app := app.NewApp()

	app.Run(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.Stdout("Service shut down with error ->", err)
		os.Exit(1)
	}
}
