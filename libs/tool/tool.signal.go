package tool

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (t *Tools) SignalHandler(callback func()) {
	usrDefChan := make(chan os.Signal, 1)
	sysSignalChan := make(chan os.Signal, 1)

	signal.Notify(usrDefChan, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for {
			sig := <-usrDefChan
			t.Stdout("User signal recv: %v", sig)
			switch sig {
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			}
		}
	}()

	signal.Notify(sysSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sysSignalChan
		t.Stdout("System signal recv: %v", sig)
		callback()
		os.Exit(0)
	}()
}

func (t *Tools) Stdout(format string, v ...any) {
	log := fmt.Sprintf("%s %s %s \n", t.GetTime(), "[go-schedule]", format)

	if _, err := fmt.Fprintf(os.Stdout, log, v...); err != nil {
		panic(err)
	}
}
