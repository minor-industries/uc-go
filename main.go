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
	rpc2 "uc-go/pkg/protocol/rpc"
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

	config, err := loadConfig(a)
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
		Logs: rpc2.NewQueue(100),
	}

	go rpc2.DecodeFrames(a.Logs, a)

	err := run(a)
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}

func loadConfig(ap *app.App) (*cfg.Config, error) {
	initConfig := cfg.DefaultConfig

	_, err := ap.Lfs.Stat(ap.ConfigFile())
	if err != nil {
		return &initConfig, storage.WriteMsgp(ap.Logs, ap.Lfs, &initConfig, ap.ConfigFile())
	}

	content, err := storage.ReadFile(ap.Lfs, ap.ConfigFile())
	if err != nil {
		return &initConfig, errors.Wrap(err, "readfile")
	}

	if len(content) == 0 {
		return &initConfig, storage.WriteMsgp(ap.Logs, ap.Lfs, &initConfig, ap.ConfigFile())
	} else {
		_, err = initConfig.UnmarshalMsg(content)
		if err != nil {
			return &initConfig, errors.Wrap(err, "unmarshal")
		}
		ap.Logs.Log("loaded configfile")
		ap.Logs.Rpc("show-config", &initConfig)
	}

	return &initConfig, nil
}
