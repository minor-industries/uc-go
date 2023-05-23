package app

import (
	"os"
	"uc-go/protocol/rpc"
	"uc-go/protocol/rpc/api"
)

func log(msg string) {
	_ = rpc.Send(
		os.Stdout,
		"log",
		&api.LogRequest{Message: msg},
	)
}
