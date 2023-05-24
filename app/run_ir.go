package app

import (
	"fmt"
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"tinygo.org/x/tinyfs/littlefs"
	"uc-go/cfg"
	"uc-go/util"
)

func HandleIR(
	logs *util.StoredLogs,
	lfs *littlefs.LFS,
	config *cfg.SyncConfig,
	msgs chan irremote.Data,
	configFileName string,
) {
	for msg := range msgs {
		line := fmt.Sprintf(
			"0x%02x, 0x%02x, 0x%02x 0x%02x",
			msg.Code,
			msg.Flags,
			msg.Command,
			msg.Address,
		)

		log(line)

		switch msg.Command {
		case 0x00: // vol-
			config.ScaleDown()
		case 0x02: // vol+
			config.ScaleUp()
		case 0x10: // 1
			config.SetAnimation("rainbow1")
		case 0x11: // 2
			config.SetAnimation("rainbow2")
		case 0x12: // 2
			config.SetAnimation("bounce")
		case 0x0E:
			if err := config.Save(logs, lfs, configFileName); err != nil {
				logs.Error(errors.Wrap(err, "save config"))
			}
		}
	}
}
