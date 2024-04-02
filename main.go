package main

import (
	"github.com/minor-industries/uc-go/app/bikelights"
	"github.com/minor-industries/uc-go/pkg/logger"
	"github.com/pkg/errors"
)

func main() {
	logs := logger.DefaultLogger

	app := bikelights.App{
		Logs: logs,
	}

	err := app.Run()
	logs.Error(errors.Wrap(err, "run exited"))
}
