package rpc

import (
	"github.com/minor-industries/uc-go/pkg/protocol/rpc/api"
	"github.com/tinylib/msgp/msgp"
	"io"
	"sync"
)

type Queue struct {
	logs chan Req
	w    io.Writer
	once sync.Once
}

func NewQueue(w io.Writer, bufSize int) *Queue {
	return &Queue{
		w:    w,
		logs: make(chan Req, bufSize),
	}
}

// TODO: move out of Queue class
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

func (sl *Queue) Start() {
	sl.once.Do(func() {
		go func() {
			for req := range sl.logs {
				Send(sl.w, req.Method, req.Body)
			}
		}()
	})
}
