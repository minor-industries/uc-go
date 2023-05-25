//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app"
	"uc-go/cfg"
	"uc-go/exe/ir"
	"uc-go/leds"
	"uc-go/storage"
	"uc-go/util"
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

	err = loadConfig(a)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go a.HandleIR(
		irMsg,
	)

	sm := leds.Setup()
	go app.RunLeds(a.Cfg, sm)

	select {}
}

func main() {
	a := &app.App{
		Logs: util.NewStoredLogs(100),
	}

	go app.DecodeFrames(a)

	err := run(a)
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}

func loadConfig(ap *app.App) error {
	initConfig := cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       0,
		Length:           5.0,
		Scale:            0.5,
		MinScale:         0.3,
		ScaleIncr:        0.02,
	}

	_, err := ap.Lfs.Stat(ap.ConfigFile())
	if err != nil {
		// TODO: this will currently fail (the first time) as WriteMsgp reads old file content
		return storage.WriteMsgp(ap.Logs, ap.Lfs, &initConfig, ap.ConfigFile())
	}

	content, err := storage.ReadFile(ap.Lfs, ap.ConfigFile())
	if err != nil {
		return errors.Wrap(err, "readfile")
	}

	if len(content) == 0 {
		return storage.WriteMsgp(ap.Logs, ap.Lfs, &initConfig, ap.ConfigFile())
	} else {
		_, err = initConfig.UnmarshalMsg(content)
		if err != nil {
			return errors.Wrap(err, "unmarshal")
		}
		ap.Logs.Log("loaded configfile")
		ap.Logs.Rpc("show-config", &initConfig)
	}

	ap.Cfg = cfg.NewSyncConfig(initConfig)

	return nil
}
