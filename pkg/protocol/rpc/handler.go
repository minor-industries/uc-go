package rpc

import (
	"github.com/minor-industries/uc-go/pkg/protocol/framing"
	"github.com/pkg/errors"
	"os"
)

type Handler interface {
	Handle(method string, body []byte) error
}

type handlerFunc struct {
	f func(method string, body []byte) error
}

func (h *handlerFunc) Handle(method string, body []byte) error {
	return h.f(method, body)
}

func HandlerFunc(f func(method string, body []byte) error) *handlerFunc {
	return &handlerFunc{f: f}
}

func DecodeFrames(
	logs *Queue,
	handler Handler,
) {
	ch := make(chan []byte, 10)

	go func() {
		framing.Decode(os.Stdin, func(msg []byte) {
			ch <- msg
		})
	}()

	for msg := range ch {
		rpcMsg := &Request{}
		_, err := rpcMsg.UnmarshalMsg(msg)
		if err != nil {
			logs.Log(errors.Wrap(err, "error: unmarshal rpc").Error())
		}

		if err := handler.Handle(rpcMsg.Method, rpcMsg.Body); err != nil {
			logs.Error(errors.Wrap(err, "handle"))
		}
	}
}
