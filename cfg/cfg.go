package cfg

import (
	"bytes"
	"github.com/pkg/errors"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/storage"
	"uc-go/util"
)

//go:generate msgp

type Config struct {
	CurrentAnimation string
	NumLeds          int
	StartIndex       int
	Length           float32

	Scale     float32
	MinScale  float32
	ScaleIncr float32
}

func (c *Config) WriteConfig(
	logs *util.StoredLogs,
	lfs *littlefs.LFS,
	filename string,
) error {
	newContent, err := c.MarshalMsg(nil)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	oldContent, err := storage.ReadFile(lfs, filename)
	if err != nil {
		return errors.Wrap(err, "readfile")
	}

	if bytes.Equal(oldContent, newContent) {
		logs.Log("content was identical, skipping write")
		return nil
	}

	err = storage.WriteFile(lfs, filename, newContent)
	if err != nil {
		return errors.Wrap(err, "writefile")
	}

	logs.Log("wrote configfile")
	return nil
}
