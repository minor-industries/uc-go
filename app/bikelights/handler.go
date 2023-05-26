package bikelights

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/pkg/protocol/rpc"
)

func (a *App) Handlers() map[string]rpc.Handler {
	return map[string]rpc.Handler{
		"dump-stored-logs": rpc.HandlerFunc(func(method string, body []byte) error {
			a.Logs.Each(func(req rpc.Req) {
				rpc.Send(os.Stdout, req.Method, req.Body)
			})
			return nil
		}),

		"get-config": rpc.HandlerFunc(func(method string, body []byte) error {
			ss := a.Cfg.SnapShot()
			if err := rpc.Send(os.Stdout, "show-config", &ss); err != nil {
				a.Logs.Error(errors.Wrap(err, "send show-config"))
			}
			return nil
		}),

		"reset-config": rpc.HandlerFunc(func(method string, body []byte) error {
			if err := a.Lfs.Remove(a.ConfigFile()); err != nil {
				a.Logs.Error(errors.Wrap(err, "remove file"))
			} else {
				a.Logs.Log("config reset")
			}
			return nil
		}),
	}
}

func (a *App) Handle(method string, body []byte) error {
	//storedLogs.Log("got rpc: " + rpcMsg.Method)

	switch method {
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
		a.Logs.Log("unknown method: " + method)
	}

	return nil
}
