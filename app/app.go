package app

import (
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/app/cfg"
	"uc-go/protocol/rpc"
	"uc-go/util"
)

const (
	configFile = "/cfg.msgp"
)

type App struct {
	Logs *rpc.Queue
	Lfs  *littlefs.LFS
	Cfg  *util.SyncConfig[cfg.Config]
}

func (a *App) ConfigFile() string {
	return configFile
}
