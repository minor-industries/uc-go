package storage

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"io"
	"machine"
	"os"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/pkg/logger"
)

var (
	blockDevice = machine.Flash
	filesystem  = littlefs.New(blockDevice)
)

func Setup(storedLogs logger.Logger) (*littlefs.LFS, error) {
	lfs := filesystem.Configure(&littlefs.Config{
		CacheSize:     512,
		LookaheadSize: 512,
		BlockCycles:   100,
	})

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

func mount(logs logger.Logger, lfs *littlefs.LFS) (err error) {
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

func WriteMsgp(
	logs logger.Logger,
	lfs *littlefs.LFS,
	msg msgp.Marshaler,
	filename string,
) error {
	newContent, err := msg.MarshalMsg(nil)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	var oldContent []byte
	_, err = lfs.Stat(filename)
	switch err {
	case nil:
		oldContent, err = ReadFile(lfs, filename)
		if err != nil {
			return errors.Wrap(err, "readfile")
		}
	}

	if bytes.Equal(oldContent, newContent) {
		logs.Log("content was identical, skipping write")
		return nil
	}

	err = WriteFile(lfs, filename, newContent)
	if err != nil {
		return errors.Wrap(err, "writefile")
	}

	logs.Log("wrote configfile")
	return nil
}
