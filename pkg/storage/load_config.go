package storage

import (
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/pkg/logger"
)

type Serializer interface {
	msgp.Marshaler
	msgp.Unmarshaler
}

func LoadConfig[T Serializer](
	lfs *littlefs.LFS,
	logs logger.Logger,
	filename string,
	forceReInit bool,
	initConfig T,
) (T, error) {
	if forceReInit {
		return initConfig, WriteMsgp(logs, lfs, initConfig, filename)
	}

	_, err := lfs.Stat(filename)
	if err != nil {
		logs.Error(errors.Wrap(err, "stat failed"))
		return initConfig, WriteMsgp(logs, lfs, initConfig, filename)
	}

	content, err := ReadFile(lfs, filename)
	if err != nil {
		return initConfig, errors.Wrap(err, "readfile")
	}

	if len(content) == 0 {
		return initConfig, WriteMsgp(logs, lfs, initConfig, filename)
	} else {
		_, err = initConfig.UnmarshalMsg(content)
		if err != nil {
			return initConfig, errors.Wrap(err, "unmarshal")
		}
		logs.Log("loaded configfile")
		logs.Rpc("show-config", initConfig)
	}

	return initConfig, nil
}
