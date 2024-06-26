//go:build qtpy_rp2040

package bbq

import (
	"github.com/minor-industries/uc-go/app/tempmon"
	"github.com/minor-industries/uc-go/pkg/spi"
	"machine"
)

var cfg = tempmon.BoardCfg{
	//Rfm: rfm69_board.PinCfg{
	//	Spi: &machine.SPI0,
	//	SpiCfg: &machine.SPIConfig{
	//		Mode: 0,
	//		SCK:  machine.A3,
	//		SDO:  machine.A4,
	//		SDI:  machine.A1,
	//	},
	//
	//	Rst:  machine.A0,
	//	Intr: machine.NoPin,
	//	Csn:  machine.A2,
	//},

	led: machine.A1,

	Tcs: []*tempmon.ThermocoupleCfg{
		{
			Name: "bbq01-meat",
			Spi: &spi.Config{
				Spi: machine.SPI0,
				Config: &machine.SPIConfig{
					Mode: 1,
					SCK:  machine.SPI0_SCK_PIN,
					SDO:  machine.SPI0_SDO_PIN,
					SDI:  machine.SPI0_SDI_PIN,
				},
				Cs: machine.A0,
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
