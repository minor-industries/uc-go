package app

import (
	"os"
	"uc-go/protocol/rpc"
	"uc-go/protocol/rpc/api"
)

func log(msg string) {
	_ = rpc.Send(&api.LogRequest{Message: msg}, os.Stdout)
}
