package app

import (
	"fmt"
	"github.com/pkg/errors"
	"tinygo.org/x/drivers/irremote"
	"uc-go/storage"
)

func (a *App) HandleIR(
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

		a.Logs.Log(line)

		switch msg.Command {
		case 0x00: // vol-
			a.Cfg.ScaleDown()
		case 0x02: // vol+
			a.Cfg.ScaleUp()
		case 0x10: // 1
			a.Cfg.SetAnimation("rainbow1")
		case 0x11: // 2
			a.Cfg.SetAnimation("rainbow2")
		case 0x12: // 2
			a.Cfg.SetAnimation("bounce")
		case 0x0E:
			ss := a.Cfg.SnapShot()
			if err := storage.WriteMsgp(a.Logs, a.Lfs, &ss, configFile); err != nil {
				a.Logs.Error(errors.Wrap(err, "save config"))
			}
		}
	}
}
