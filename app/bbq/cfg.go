package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin
}

var cfg = BoardCfg{
	Rfm: rfm69_board.PinCfg{
		Spi: machine.SPI0,

		Rst:  machine.GPIO6,
		Intr: machine.GPIO7,

		Sck: machine.GPIO2,
		Sdo: machine.GPIO3,
		Sdi: machine.GPIO4,
		Csn: machine.GPIO5,
	},

	led: machine.LED,
}
