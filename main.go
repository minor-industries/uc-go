package main

import (
	"fmt"
	"github.com/pkg/errors"
	"uc-go/app/bbq"
)

type logger struct{}

func (l *logger) Log(s string) {
	fmt.Println(s)
}

func (l *logger) Error(err error) {
	fmt.Println("error: " + err.Error())
}

func (l *logger) Rpc(s string, i interface{}) error {
	fmt.Println("rpc: " + s)
	return nil
}

func main() {
	logs := &logger{}
	err := bbq.Run(logs)
	logs.Error(errors.Wrap(err, "bbq exited"))
}
