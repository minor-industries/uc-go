package app

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/protocol/framing"
	"uc-go/protocol/rpc"
)

func DecodeFrames(a *App) {
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
			a.Logs.Log(errors.Wrap(err, "error: unmarshal rpc").Error())
		}

		//storedLogs.Log("got rpc: " + rpcMsg.Method)

		switch rpcMsg.Method {
		case "dump-stored-logs":
			a.Logs.Each(func(req rpc.Req) {
				rpc.Send(os.Stdout, req.Method, req.Body)
			})

		case "get-config":
			ss := a.Cfg.SnapShot()
			if err := rpc.Send(os.Stdout, "show-config", &ss); err != nil {
				a.Logs.Error(errors.Wrap(err, "send show-config"))
			}

		case "reset-config":
			if err := a.Lfs.Remove(a.ConfigFile()); err != nil {
				a.Logs.Error(errors.Wrap(err, "remove file"))
			} else {
				a.Logs.Log("config reset")
			}

		default:
			a.Logs.Log("unknown method: " + rpcMsg.Method)
		}
	}
}
