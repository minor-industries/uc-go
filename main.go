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

	go rpc.DecodeFrames(a.Logs, a)

	err := a.Run()
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}
