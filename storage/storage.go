package storage

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"machine"
	"os"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/util"
)

var (
	blockDevice = machine.Flash
	filesystem  = littlefs.New(blockDevice)
)

func Setup(storedLogs *util.StoredLogs) (*littlefs.LFS, error) {
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

	if err := mount(storedLogs, lfs); err != nil {
		return nil, errors.Wrap(err, "mount")
	}

	storedLogs.Log("mounted")

	n, err := lfs.Size()
	if err != nil {
		return nil, errors.Wrap(err, "size")
	}

	storedLogs.Log(fmt.Sprintf("size = %d", n))

	root, err := lfs.Open("/")
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}

	infos, err := root.Readdir(0)
	if err != nil {
		return nil, errors.Wrap(err, "readdir")
	}

	for _, info := range infos {
		storedLogs.Log(fmt.Sprintf("file: %s %d", info.Name(), info.Size()))
	}

	return lfs, nil
}

func mount(logs *util.StoredLogs, lfs *littlefs.LFS) (err error) {
	for i := 0; i <= 1; i++ {
		err = lfs.Mount()
		if err != nil {
			if err := lfs.Format(); err != nil {
				return errors.Wrap(err, "format")
			}
			logs.Log("formatted")
			continue
		}
	}
	return
}

func ReadFile(
	lfs *littlefs.LFS,
	name string,
) ([]byte, error) {
	fp, err := lfs.Open(name)
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	defer fp.Close()

	content, err := io.ReadAll(fp)
	if err != nil {
		return nil, errors.Wrap(err, "readall")
	}

	return content, nil
}

func WriteFile(
	lfs *littlefs.LFS,
	name string,
	content []byte,
) error {
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
