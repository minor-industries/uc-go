package app

import (
	"fmt"
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
		reply := fmt.Sprintf("got frame: [%s]", msg)
		log(reply)

		rpcMsg := &rpc.Request{}
		_, err := rpcMsg.UnmarshalMsg(msg)
		if err != nil {
			storedLogs.Log(errors.Wrap(err, "error: unmarshal rpc").Error())
		}

		storedLogs.Log("got rpc: " + rpcMsg.Method)

		switch rpcMsg.Method {
		case "dump-stored-logs":
			storedLogs.Each(func(s string) {
				log(s)
			})

		default:
			storedLogs.Log("unknown method: " + rpcMsg.Method)
		}
	}
}
