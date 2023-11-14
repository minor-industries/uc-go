//go:build pico
// +build pico

package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
)

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

	Tcs: []*ThermocoupleCfg{
		//{
		//	Name: "bbq01-meat",
		//	Spi: &spi.Config{
		//		Spi: machine.SPI0,
		//		Config: &machine.SPIConfig{
		//			Mode: 1,
		//			SCK:  machine.GPIO2,
		//			SDO:  machine.GPIO3,
		//			SDI:  machine.GPIO4,
		//		},
		//		Cs: machine.GPIO8,
		//	},
		//},
		//{
		//	Name: "bbq01-bbq",
		//	Spi: &spi.Config{
		//		Spi: machine.SPI0,
		//		Config: &machine.SPIConfig{
		//			Mode: 1,
		//			SCK:  machine.GPIO2,
		//			SDO:  machine.GPIO3,
		//			SDI:  machine.GPIO4,
		//		},
		//		Cs: machine.GPIO9,
		//	},
		//},
	},
}
