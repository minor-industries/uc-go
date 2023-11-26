package tempmon

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin

	i2cCfg *machine.I2CConfig
	i2c    *machine.I2C
}
