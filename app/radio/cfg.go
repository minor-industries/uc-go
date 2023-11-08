package radio

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	i2c *machine.I2C
	sda machine.Pin
	scl machine.Pin

	// misc
	led machine.Pin
}
