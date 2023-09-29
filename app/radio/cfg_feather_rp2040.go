//go:build feather_rp2040

package radio

import "machine"

var cfg = pinCfg{
	spi:  machine.SPI1,
	rst:  machine.GPIO17,
	intr: machine.GPIO21,
	sck:  machine.GPIO14,
	sdo:  machine.GPIO15,
	sdi:  machine.GPIO8,
	csn:  machine.GPIO16,

	i2c: machine.I2C1,
	sda: machine.GPIO2,
	scl: machine.GPIO3,

	led: machine.GPIO13,
}
