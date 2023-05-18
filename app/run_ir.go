package app

import (
	"fmt"
	"tinygo.org/x/drivers/irremote"
	"uc-go/cfg"
)

func HandleIR(
	config *cfg.SyncConfig,
	msgs chan irremote.Data,
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
		}
	}
}