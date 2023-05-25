package app

import (
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

func (a *App) ConfigFile() string {
	return configFile
}
