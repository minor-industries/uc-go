package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

type TCCfg struct {
	Spi *machine.SPI
	Csn machine.Pin
}

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin

	Tc TCCfg
}

var cfg = BoardCfg{
	Rfm: rfm69_board.PinCfg{
		Spi: machine.SPI0,
		SpiCfg: &machine.SPIConfig{
			Mode: 0,
			SCK:  machine.GPIO2,
			SDO:  machine.GPIO3,
			SDI:  machine.GPIO4,
		},

		Rst:  machine.GPIO6,
		Intr: machine.GPIO7,
		Csn:  machine.GPIO5,
	},

	led: machine.LED,

	Tc: TCCfg{
		Spi: machine.SPI0,
		Csn: machine.GPIO8,
	},
}
