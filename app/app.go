package app

import (
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/cfg"
	"uc-go/util"
)

const (
	configFile = "/cfg.msgp"
)

type App struct {
	Logs *util.StoredLogs
	Lfs  *littlefs.LFS
	Cfg  *cfg.SyncConfig
}

func (a *App) ConfigFile() string {
	return configFile
}
