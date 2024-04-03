//go:build itsybitsy_m4

package leds

import (
	"device/sam"
	neopixel_spi "github.com/minor-industries/neopixel-spi"
	"github.com/minor-industries/neopixel-spi/default_driver"
	"github.com/pkg/errors"
	"image/color"
	"machine"
)

var TxFullCounter int64 // TODO

func Setup(numLeds int) (*neopixel_spi.NeoSpiDriver, error) {
	d := default_driver.Configure(&neopixel_spi.Cfg{
		SPI:        &machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5},
		SCK:        machine.PA22, // 5.1 (sercom alt)
		SDO:        machine.PA23, // 5.0 (sercom alt)
		SDI:        machine.PA20, // 5.2 (sercom alt)
		LedCount:   numLeds,
		SpaceCount: 2000,
	})

	if err := d.Init(); err != nil {
		return nil, errors.Wrap(err, "inii")
	}

	return d, nil
}

func Write(driver_ any, pixels []color.RGBA) {
	d := driver_.(*neopixel_spi.NeoSpiDriver)
	d.Animate(pixels)
}
