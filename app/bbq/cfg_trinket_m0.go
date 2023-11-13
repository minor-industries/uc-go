//go:build trinket_m0
// +build trinket_m0

package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

var cfg = BoardCfg{
	Rfm: rfm69_board.PinCfg{
		Spi: &machine.SPI0,
		SpiCfg: &machine.SPIConfig{
			Mode: 0,
			SCK:  machine.A3,
			SDO:  machine.A4,
			SDI:  machine.A1,
		},

		Rst:  machine.NoPin,
		Intr: machine.NoPin,
		Csn:  machine.A2,
	},

	led: machine.LED,

	Tc0: &ThermocoupleCfg{
		Name: "bbq01-meat",
		Spi: &spi.Config{
			Spi: &machine.SPI0,
			Config: &machine.SPIConfig{
				Mode: 1,
				SCK:  machine.A3,
				SDO:  machine.A4,
				SDI:  machine.A1,
			},
			Cs: machine.A0,
		},
	},

	Tc1: &ThermocoupleCfg{
		Name: "bbq01-bbq",
		Spi: &spi.Config{
			Spi: &machine.SPI0,
			Config: &machine.SPIConfig{
				Mode: 1,
				SCK:  machine.A3,
				SDO:  machine.A4,
				SDI:  machine.A1,
			},
			Cs: machine.A0, // TODO:
		},
	},
}
