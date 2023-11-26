//go:build feather_rp2040

package tempmon

import "machine"

// feather with rfm69 built-in
var cfg1 = PinCfg{
	spi: machine.SPI1,

	rst:  machine.GPIO17,
	intr: machine.GPIO21,

	sck: machine.GPIO14,
	sdo: machine.GPIO15,
	sdi: machine.GPIO8,
	csn: machine.GPIO16,

	i2c: machine.I2C1,
	sda: machine.GPIO2,
	scl: machine.GPIO3,

	led: machine.GPIO13,
}

// feather with rfm69 hat
var cfg2 = PinCfg{
	spi: machine.SPI0,

	rst:  machine.GPIO11,
	intr: machine.GPIO21,

	sck: machine.SPI0_SCK_PIN,
	sdo: machine.SPI0_SDO_PIN,
	sdi: machine.SPI0_SDI_PIN,
	csn: machine.GPIO9,

	i2c: machine.I2C1,
	sda: machine.GPIO2,
	scl: machine.GPIO3,

	led: machine.GPIO13,
}

var cfg = cfg2
