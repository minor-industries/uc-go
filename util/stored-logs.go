package util

import (
	"github.com/tinylib/msgp/msgp"
	"uc-go/protocol/rpc"
	"uc-go/protocol/rpc/api"
)

type StoredLogs struct {
	logs chan rpc.Req
}

func NewStoredLogs(bufSize int) *StoredLogs {
	return &StoredLogs{logs: make(chan rpc.Req, bufSize)}
}

func (sl *StoredLogs) Log(s string) {
	sl.Rpc("log", api.LogRequest{Message: s})
}

func (sl *StoredLogs) Rpc(method string, msg msgp.Marshaler) {
	// drop logs if buffer is full
	select {
	case sl.logs <- rpc.Req{
		Method: method,
		Body:   msg,
	}:
	default:
	}
}

func (sl *StoredLogs) Error(err error) {
	sl.Log("error: " + err.Error())
}

func (sl *StoredLogs) Each(cb func(rpc.Req)) {
	for {
		select {
		case s := <-sl.logs:
			cb(s)
		default:
			return
		}
	}
}
