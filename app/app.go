package app

import (
	"github.com/pkg/errors"
	"os"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/cfg"
	"uc-go/protocol/rpc"
)

const (
	configFile = "/cfg.msgp"
)

type App struct {
	Logs *rpc.Queue
	Lfs  *littlefs.LFS
	Cfg  *cfg.SyncConfig
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

func (a *App) ConfigFile() string {
	return configFile
}
