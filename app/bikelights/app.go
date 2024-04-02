package bikelights

import (
	"fmt"
	"github.com/minor-industries/uc-go/app/bikelights/cfg"
	"github.com/minor-industries/uc-go/pkg/ir"
	"github.com/minor-industries/uc-go/pkg/leds"
	"github.com/minor-industries/uc-go/pkg/logger"
	"github.com/minor-industries/uc-go/pkg/storage"
	"github.com/minor-industries/uc-go/pkg/util"
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"tinygo.org/x/tinyfs/littlefs"
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
	//<-time.After(5 * time.Second)
	a.Logs.Log("Hello")

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

	fmt.Println("here 1")

	a.Cfg = util.NewSyncConfig(*config)

	fmt.Println("here 2")

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	fmt.Println("here 3")

	go a.HandleIR(
		irMsg,
	)

	fmt.Println("here 4")

	sm, err := leds.Setup(a.Cfg.SnapShot().NumLeds)
	if err != nil {
		return errors.Wrap(err, "setup leds")
	}

	fmt.Println("here 5")

	runLeds(a.Cfg, sm)

	return errors.New("unexpected exit")
}
