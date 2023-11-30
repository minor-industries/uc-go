//go:build qtpy

package tempmon

import (
	"image/color"
	"machine"
	"tinygo.org/x/drivers/ws2812"
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

type neoBlinker struct {
	pixel *ws2812.Device
}

func (n *neoBlinker) Set(on bool) {
	// TODO: don't hardcode these colors, give a function to change them
	if on {
		n.pixel.WriteColors([]color.RGBA{{0, 0, 16, 0}})
	} else {
		n.pixel.WriteColors([]color.RGBA{{0, 0, 0, 0}})
	}
}

var blinker neoBlinker

func init() {
	machine.NEOPIXELS_POWER.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXELS_POWER.High()

	machine.NEOPIXELS.Configure(machine.PinConfig{Mode: machine.PinOutput})

	device := ws2812.New(machine.NEOPIXELS)
	blinker = neoBlinker{pixel: &device}
}
