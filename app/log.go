package app

import (
	"os"
	"uc-go/protocol/framing"
	"uc-go/protocol/rpc"
	"uc-go/protocol/rpc/api"
)

func log(msg string) {
	req := &api.LogRequest{Message: msg}
	marshal, err := req.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	rpcMsg := rpc.Request{
		Method: "log",
		Body:   marshal,
	}

	rpcMarshal, err := rpcMsg.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	frame := framing.Encode(rpcMarshal)
	_, _ = os.Stdout.Write(frame)
}
