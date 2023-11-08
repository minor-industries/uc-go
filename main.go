//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/app/bbq"
	"uc-go/pkg/protocol/rpc"
)

func main() {
	logs := rpc.NewQueue(os.Stdout, 100)
	router := rpc.NewRouter()
	_ = router.Register(map[string]rpc.Handler{
		"__sys__.dump-stored-logs": rpc.HandlerFunc(func(method string, body []byte) error {
			logs.Start()
			return nil
		}),
	})
	go rpc.DecodeFrames(logs, router)

	if err := bbq.Run(logs); err != nil {
		logs.Error(errors.Wrap(err, "run"))
	}

	select {}
}
