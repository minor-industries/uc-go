package app

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/protocol/framing"
	"uc-go/protocol/rpc"
	"uc-go/util"
)

func DecodeFrames(storedLogs *util.StoredLogs) {
	ch := make(chan []byte, 10)

	go func() {
		framing.Decode(os.Stdin, func(msg []byte) {
			ch <- msg
		})
	}()

	for msg := range ch {
		rpcMsg := &rpc.Request{}
		_, err := rpcMsg.UnmarshalMsg(msg)
		if err != nil {
			storedLogs.Log(errors.Wrap(err, "error: unmarshal rpc").Error())
		}

		//storedLogs.Log("got rpc: " + rpcMsg.Method)

		switch rpcMsg.Method {
		case "dump-stored-logs":
			storedLogs.Each(func(s string) {
				log("stored: " + s)
			})

		case "get-config":
			rpc.Send(os.Stdout, "show-config", nil)

		default:
			storedLogs.Log("unknown method: " + rpcMsg.Method)
		}
	}
}
