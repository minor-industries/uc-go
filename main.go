//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/app/bikelights"
	"uc-go/pkg/protocol/rpc"
)

func main() {
	rpcQueue := rpc.NewQueue(os.Stdout, 100)
	a := &bikelights.App{
		Logs: rpcQueue,
	}

	router := rpc.NewRouter()
	router.Register(map[string]rpc.Handler{
		"__sys__.dump-stored-logs": rpc.HandlerFunc(func(method string, body []byte) error {
			rpcQueue.Start()
			return nil
		}),
	})

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
