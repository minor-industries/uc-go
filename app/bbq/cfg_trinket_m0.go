//go:build trinket_m0
// +build trinket_m0

package bbq

import (
	"machine"
	"uc-go/app/tempmon"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

var cfg = tempmon.BoardCfg{
	Rfm: rfm69_board.PinCfg{
		Spi: &spi.Config{
			Spi: &machine.SPI0,
			Config: &machine.SPIConfig{
				Mode: 0,
				SCK:  machine.SPI0_SCK_PIN,
				SDO:  machine.SPI0_SDO_PIN,
				SDI:  machine.SPI0_SDI_PIN,
			},
			Cs: machine.D1,
		},
		Rst:  machine.NoPin,
		Intr: machine.NoPin,
	},

	led: machine.LED,

	Tcs: []*tempmon.ThermocoupleCfg{
		{
			Name: "bbq01-meat",
			Spi: &spi.Config{
				Spi: &machine.SPI0,
				Config: &machine.SPIConfig{
					Mode: 1,
					SCK:  machine.SPI0_SCK_PIN,
					SDO:  machine.SPI0_SDO_PIN,
					SDI:  machine.SPI0_SDI_PIN,
				},
				Cs: machine.D0,
			},
		},
		//{
		//	Name: "bbq01-bbq",
		//	Spi: &spi.Config{
		//		Spi: &machine.SPI0,
		//		Config: &machine.SPIConfig{
		//			Mode: 1,
		//			SCK:  machine.A3,
		//			SDO:  machine.A4,
		//			SDI:  machine.A1,
		//		},
		//		Cs: machine.A0, // TODO:
		//	},
		//},
	},
}
