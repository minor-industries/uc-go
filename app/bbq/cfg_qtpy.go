//go:build qtpy

package bbq

import (
	"github.com/minor-industries/uc-go/app/tempmon"
	rfm69_board "github.com/minor-industries/uc-go/pkg/rfm69-board"
	"github.com/minor-industries/uc-go/pkg/spi"
	"machine"
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
			Cs: machine.A1,
		},
		Rst:  machine.NoPin,
		Intr: machine.NoPin,
	},

	led: machine.NoPin,
}
