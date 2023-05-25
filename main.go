//go:build rp2040

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app"
	"uc-go/app/cfg"
	"uc-go/pkg/ir"
	"uc-go/pkg/leds"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
	"uc-go/wifi"
)

func run(a *app.App) error {
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
	go app.RunLeds(a.Cfg, sm)

	r := wifi.F(2, 3)
	a.Logs.Log(fmt.Sprintf("F(2,3) = %d", r))

	select {}
}

func main() {
	a := &app.App{
		Logs: rpc.NewQueue(100),
	}

	go rpc.DecodeFrames(a.Logs, a)

	err := run(a)
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}
