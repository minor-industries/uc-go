package rpc

import (
	"github.com/tinylib/msgp/msgp"
	"uc-go/protocol/rpc/api"
)

type Queue struct {
	logs chan Req
}

func NewQueue(bufSize int) *Queue {
	return &Queue{logs: make(chan Req, bufSize)}
}

func (sl *Queue) Log(s string) {
	sl.Rpc("log", api.LogRequest{Message: s})
}

func (sl *Queue) Rpc(method string, msg msgp.Marshaler) {
	// drop logs if buffer is full
	select {
	case sl.logs <- Req{
		Method: method,
		Body:   msg,
	}:
	default:
	}
}

func (sl *Queue) Error(err error) {
	sl.Log("error: " + err.Error())
}

func (sl *Queue) Each(cb func(Req)) {
	for {
		select {
		case s := <-sl.logs:
			cb(s)
		default:
			return
		}
	}
}
