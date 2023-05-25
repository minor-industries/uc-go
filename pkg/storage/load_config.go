package storage

import (
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/pkg/protocol/rpc"
)

type Serializer interface {
	msgp.Marshaler
	msgp.Unmarshaler
}

func LoadConfig[T Serializer](
	lfs *littlefs.LFS,
	logs *rpc.Queue,
	filename string,
	initConfig T,
) (T, error) {
	_, err := lfs.Stat(filename)
	if err != nil {
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
