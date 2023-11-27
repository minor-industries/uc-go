package setaddr

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
	"uc-go/pkg/logger"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/storage"
)

const newAddr = 0xD0

func Run(logs logger.Logger) error {
	<-time.After(5 * time.Second)

	env, err := rfm69_board.LoadConfig(logs)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	ss := env.SnapShot()
	fmt.Printf("address is: 0x%02x\n", ss.NodeAddr)

	if ss.NodeAddr != newAddr {
		ss.NodeAddr = newAddr

		lfs, err := storage.Setup(logs)
		if err != nil {
			return errors.Wrap(err, "setup storage")
		}

		if lfs == nil {
			return errors.New("no lfs")
		}

		if err := storage.WriteMsgp(logs, lfs, &ss, "/radio-cfg.msgp"); err != nil {
			return errors.Wrap(err, "write msgp")
		}

		fmt.Printf("set new address to 0x%02x\n", ss.NodeAddr)
	}

	return nil
}
