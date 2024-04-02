package tempmon

import (
	rfm69_board "github.com/minor-industries/uc-go/pkg/rfm69-board"
	"machine"
)

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin

	i2cCfg *machine.I2CConfig
	i2c    *machine.I2C
}
