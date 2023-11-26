//go:build qtpy

package tempmon

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

var cfg = BoardCfg{
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

	i2cCfg: &machine.I2CConfig{},
	i2c:    machine.I2C0,
}
