//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/app/bikelights"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/storage"
)

func main2() {
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

func resetConfig() {
	logs := rpc.NewQueue(os.Stdout, 100)

	lfs, err := storage.Setup(logs)
	if err != nil {
		panic("no")
	}

	if err := storage.WriteMsgp(logs, lfs, &cfg.DefaultConfig, "/cfg.msgp"); err != nil {
		panic("no")
	}

	select {}
}

func setConfig() {
	logs := rpc.NewQueue(os.Stdout, 100)

	lfs, err := storage.Setup(logs)
	if err != nil {
		panic("no")
	}

	c := cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       10,
		Length:           5,
		Scale:            0.5,
		MinScale:         0.04,
		ScaleIncr:        0.02,
	}

	if err := storage.WriteMsgp(logs, lfs, &c, "/cfg.msgp"); err != nil {
		panic("no")
	}
}

func main() {
	main2()
	//resetConfig()
	//setConfig()
}
