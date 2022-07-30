package main

import (
	"./app"
	"./param"
	"./serverError"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

func cleanupOnInterrupt(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT)

	go func() {
		<-chSignal
		appInst.Close()
		os.Exit(0)
	}()
}

func reOpenLogOnHup(appInst *app.App) {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGHUP)

	go func() {
		for range chSignal {
			errs := appInst.ReOpenLog()
			serverError.CheckFatal(errs...)
		}
	}()
}

func main() {
	params := param.ParseCli()
	appInst, errs := app.NewApp(params)
	serverError.CheckFatal(errs...)

	if appInst == nil {
		serverError.CheckFatal(errors.New("failed to create application instance"))
	}

	cleanupOnInterrupt(appInst)
	reOpenLogOnHup(appInst)
	errs = appInst.Open()
	serverError.CheckFatal(errs...)

	appInst.Close()
}
