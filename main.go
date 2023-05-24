//go:build rp2040

package main

import (
	"github.com/pkg/errors"
	"io"
	"os"
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

	loadedCfg := &cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       0,
		Length:           5.0,
		Scale:            0.5,
		MinScale:         0.3,
		ScaleIncr:        0.02,
	}

	err = loadConfig(storedLogs, lfs, loadedCfg)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	config := cfg.NewSyncConfig(*loadedCfg)

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

func main() {
	storedLogs := util.NewStoredLogs(100)

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

func readFile(lfs *littlefs.LFS, name string) ([]byte, error) {
	fp, err := lfs.Open("/cfg.msgp")
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	defer fp.Close()

	content, err := io.ReadAll(fp)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}

	return content, nil
}

func writeFile(lfs *littlefs.LFS, name string, content []byte) error {
	fp, err := lfs.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return errors.Wrap(err, "openfile")
	}
	defer fp.Close()

	_, err = fp.Write(content)
	if err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func loadConfig(logs *util.StoredLogs, lfs *littlefs.LFS, c *cfg.Config) error {
	content, err := readFile(lfs, configFile)
	if err != nil {
		return errors.Wrap(err, "read all")
	}

	if len(content) == 0 {
		newContent, err := c.MarshalMsg(nil)
		if err != nil {
			return errors.Wrap(err, "marshal")
		}

		err = writeFile(lfs, configFile, newContent)
		if err != nil {
			return errors.Wrap(err, "writefile")
		}

		logs.Log("wrote configfile")
	} else {
		_, err = c.UnmarshalMsg(content)
		if err != nil {
			return errors.Wrap(err, "unmarshal")
		}

		logs.Log("loaded configfile")
	}

	return nil
}
