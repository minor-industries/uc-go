//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"uc-go/app/bikelights"
	"uc-go/pkg/protocol/rpc"
)

func main() {
	a := &bikelights.App{
		Logs: rpc.NewQueue(100),
	}

	router := rpc.NewRouter()
	router.Register(a.Handlers())

	go rpc.DecodeFrames(a.Logs, router)

	err := a.Run()
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}
