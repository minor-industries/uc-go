//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"io"
	"tinygo.org/x/drivers/irremote"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/app"
	"uc-go/cfg"
	"uc-go/exe/ir"
	"uc-go/leds"
	"uc-go/storage"
	"uc-go/util"
)

func main() {
	storedLogs := util.NewStoredLogs(100)
	lfs, err := storage.Setup(storedLogs)
	if err != nil {
		storedLogs.Error(errors.Wrap(err, "setup storage"))
	}

	var config *cfg.SyncConfig
	{
		c := &cfg.Config{
			CurrentAnimation: "rainbow1",
			NumLeds:          150,
			StartIndex:       0,
			Length:           5.0,
			Scale:            0.5,
			MinScale:         0.3,
			ScaleIncr:        0.02,
		}

		err = loadConfig(lfs, c)
		if err != nil {
			storedLogs.Error(errors.Wrap(err, "load config"))
		}

		config = cfg.NewSyncConfig(*c)
	}

	go app.DecodeFrames(storedLogs)

	irMsg := make(chan irremote.Data, 10)
	ir.Main(func(data irremote.Data) {
		irMsg <- data
	})

	go app.HandleIR(config, irMsg)

	sm := leds.Setup()
	go app.RunLeds(config, sm)

	select {}
}

func loadConfig(lfs *littlefs.LFS, c *cfg.Config) error {
	fp, err := lfs.Open("/cfg.msgp")
	if err != nil {
		return errors.Wrap(err, "open")
	}

	content, err := io.ReadAll(fp)
	if err != nil {
		return errors.Wrap(err, "read all")
	}

	_, err = c.UnmarshalMsg(content)
	if err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	return nil
}
