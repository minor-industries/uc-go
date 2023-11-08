//go:build pico

package radio

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

var cfg = BoardCfg{
	Rfm: rfm69_board.PinCfg{
		Spi: machine.SPI0,

		Rst:  machine.GP6,
		Intr: machine.GP7,

		Sck: machine.GP2,
		Sdo: machine.GP3,
		Sdi: machine.GP4,
		Csn: machine.GP5,
	},
	i2c: machine.I2C0,
	sda: machine.GP0,
	scl: machine.GP1,

	led: machine.LED,
}
