//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/app"
	"uc-go/cfg"
	"uc-go/exe/ir"
	"uc-go/leds"
	"uc-go/storage"
	"uc-go/util"
)

func run(storedLogs *util.StoredLogs) error {
	lfs, err := storage.Setup(storedLogs)
	if err != nil {
		return errors.Wrap(err, "setup storage")
	}

	if lfs == nil {
		return errors.New("no lfs")
	}

	loadedCfg := &cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       0,
		Length:           5.0,
		Scale:            0.5,
		MinScale:         0.3,
		ScaleIncr:        0.02,
	}

	//err = loadConfig(storedLogs, lfs, loadedCfg)
	//if err != nil {
	//	return errors.Wrap(err, "load config")
	//}

	config := cfg.NewSyncConfig(*loadedCfg)

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go app.HandleIR(
		storedLogs,
		lfs,
		config,
		irMsg,
		configFile,
	)

	sm := leds.Setup()
	go app.RunLeds(config, sm)

	select {}
}

func main() {
	storedLogs := util.NewStoredLogs(100)
	go app.DecodeFrames(storedLogs)

	err := run(storedLogs)
	if err != nil {
		storedLogs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		storedLogs.Log("run exited")
	}

	select {}
}

const (
	configFile = "/cfg.msgp"
)

func loadConfig(logs *util.StoredLogs, lfs *littlefs.LFS, c *cfg.Config) error {
	_, err := lfs.Stat(configFile)
	if err != nil {
		// TODO: this will currently fail (the first time) as WriteMsgp reads old file content
		return storage.WriteMsgp(logs, lfs, c, configFile)
	}

	content, err := storage.ReadFile(lfs, configFile)
	if err != nil {
		return errors.Wrap(err, "readfile")
	}

	if len(content) == 0 {
		return storage.WriteMsgp(logs, lfs, c, configFile)
	} else {
		_, err = c.UnmarshalMsg(content)
		if err != nil {
			return errors.Wrap(err, "unmarshal")
		}

		logs.Log("loaded configfile")
	}

	return nil
}
