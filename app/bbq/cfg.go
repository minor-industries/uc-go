package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin

	Tc *spi.Config
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
		Csn:  machine.GPIO9,
	},

	led: machine.LED,

	Tc: &spi.Config{
		Spi: machine.SPI0,
		Config: &machine.SPIConfig{
			Mode: 1,
			SCK:  machine.GPIO2,
			SDO:  machine.GPIO3,
			SDI:  machine.GPIO4,
		},
		Cs: machine.GPIO8,
	},
}
