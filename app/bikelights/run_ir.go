package bikelights

import (
	"fmt"
	"github.com/minor-industries/uc-go/app/bikelights/cfg"
	"github.com/minor-industries/uc-go/pkg/protocol/rpc"
	"github.com/minor-industries/uc-go/pkg/storage"
	"github.com/minor-industries/uc-go/pkg/util"
	"github.com/pkg/errors"
	"os"
	"tinygo.org/x/drivers/irremote"
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
			a.Cfg.Edit(func(c *cfg.Config) {
				c.Scale = util.Clamp(c.MinScale, c.Scale-c.ScaleIncr, 1.0)
			})
		case 0x02: // vol+
			a.Cfg.Edit(func(c *cfg.Config) {
				c.Scale = util.Clamp(c.MinScale, c.Scale+c.ScaleIncr, 1.0)
			})
		case 0x10: // 1
			a.Cfg.Edit(func(c *cfg.Config) {
				c.CurrentAnimation = "rainbow1"
			})
		case 0x11: // 2
			a.Cfg.Edit(func(c *cfg.Config) {
				c.CurrentAnimation = "rainbow2"
			})
		case 0x12: // 2
			a.Cfg.Edit(func(c *cfg.Config) {
				c.CurrentAnimation = "bounce"
			})
		case 0x19: // 8
			a.Cfg.Edit(func(c *cfg.Config) {
				c.CurrentAnimation = "white"
			})
		case 0x0E: // the "return" button
			ss := a.Cfg.SnapShot()
			if err := storage.WriteMsgp(a.Logs, a.Lfs, &ss, configFile); err != nil {
				a.Logs.Error(errors.Wrap(err, "save config"))
			}
		case 0x04: // setup
			ss := a.Cfg.SnapShot()
			_ = rpc.Send(os.Stdout, "show-config", &ss)
		}
	}
}
