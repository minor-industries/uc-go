//go:build pico
// +build pico

package bbq

import (
	"github.com/minor-industries/uc-go/app/tempmon"
	"github.com/minor-industries/uc-go/pkg/spi"
	"machine"
)

var cfg = tempmon.BoardCfg{
	//Rfm: rfm69_board.PinCfg{
	//	Spi: machine.SPI0,
	//	SpiCfg: &machine.SPIConfig{
	//		Mode: 0,
	//		SCK:  machine.GPIO2,
	//		SDO:  machine.GPIO3,
	//		SDI:  machine.GPIO4,
	//	},
	//
	//	Rst:  machine.GPIO6,
	//	Intr: machine.GPIO7,
	//	Csn:  machine.GPIO5,
	//},

	led: machine.LED,

	Tcs: []*tempmon.ThermocoupleCfg{
		{
			Name: "bbq01-meat",
			Spi: &spi.Config{
				Spi: machine.SPI0,
				Config: &machine.SPIConfig{
					Mode: 1,
					SCK:  machine.GPIO2,
					SDO:  machine.GPIO3,
					SDI:  machine.GPIO4,
				},
				Cs: machine.GPIO5,
			},
		},
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
