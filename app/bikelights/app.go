package bikelights

import (
	"fmt"
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/ir"
	"uc-go/pkg/leds"
	"uc-go/pkg/logger"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
	"uc-go/wifi"
)

const (
	configFile = "/cfg.msgp"
)

type App struct {
	Logs logger.Logger
	Lfs  *littlefs.LFS
	Cfg  *util.SyncConfig[cfg.Config]
}

func (a *App) ConfigFile() string {
	return configFile
}

func (a *App) Run() error {
	var err error

	a.Lfs, err = storage.Setup(a.Logs)
	if err != nil {
		return errors.Wrap(err, "setup storage")
	}

	if a.Lfs == nil {
		return errors.New("no lfs")
	}

	config, err := storage.LoadConfig[*cfg.Config](
		a.Lfs,
		a.Logs,
		a.ConfigFile(),
		false,
		&cfg.DefaultConfig,
	)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	a.Cfg = util.NewSyncConfig(*config)

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go a.HandleIR(
		irMsg,
	)

	sm := leds.Setup()
	go runLeds(a.Cfg, sm)

	r := wifi.F(2, 3)
	a.Logs.Log(fmt.Sprintf("F(2,3) = %d", r))

	select {}
}
