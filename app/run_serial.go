package app

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/protocol/framing"
	"uc-go/protocol/rpc"
)

type Handler interface {
	Handle(string, []byte) error
}

func DecodeFrames(
	logs *rpc.Queue,
	handler Handler,
) {
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
			logs.Log(errors.Wrap(err, "error: unmarshal rpc").Error())
		}

		if err := handler.Handle(rpcMsg.Method, rpcMsg.Body); err != nil {
			logs.Error(errors.Wrap(err, "handle"))
		}
	}
}
