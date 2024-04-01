package main

import (
	"github.com/pkg/errors"
	"uc-go/app/bikelights"
	"uc-go/pkg/logger"
)

func main() {
	logs := logger.DefaultLogger

	app := bikelights.App{
		Logs: logs,
	}

	err := app.Run()
	logs.Error(errors.Wrap(err, "run exited"))
}
