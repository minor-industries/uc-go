package storage

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"os"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/util"
)

var (
	blockDevice = machine.Flash
	filesystem  = littlefs.New(blockDevice)
)

func Setup(storedLogs *util.StoredLogs) error {
	lfs := filesystem.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	storedLogs.Log(fmt.Sprintf(
		"lsblk start=0x%x, end=0x%x",
		machine.FlashDataStart(),
		machine.FlashDataEnd(),
	))

	err := lfs.Format()
	if err != nil {
		return errors.Wrap(err, "format")
	}

	err = lfs.Mount()
	if err != nil {
		return errors.Wrap(err, "mount")
	}

	storedLogs.Log("mounted")

	n, err := lfs.Size()
	if err != nil {
		return errors.Wrap(err, "size")
	}

	storedLogs.Log(fmt.Sprintf("size = %d", n))

	file, err := lfs.OpenFile("/cfg.msgp", os.O_CREATE)
	if err != nil {
		return errors.Wrap(err, "create")
	}

	storedLogs.Log(fmt.Sprintf("file= %v", file))

	root, err := lfs.Open("/")
	if err != nil {
		return errors.Wrap(err, "open")
	}

	infos, err := root.Readdir(0)
	if err != nil {
		return errors.Wrap(err, "readdir")
	}

	for _, info := range infos {
		storedLogs.Log(fmt.Sprintf("file: %s %d", info.Name(), info.Size()))
	}

	return nil
}
