//go:build pico

package radio

import "machine"

var cfg = pinCfg{
	spi: machine.SPI0,

	rst:  machine.GP6,
	intr: machine.GP7,

	sck: machine.GP2,
	sdo: machine.GP3,
	sdi: machine.GP4,
	csn: machine.GP5,

	i2c: machine.I2C0,
	sda: machine.GP0,
	scl: machine.GP1,

	led: machine.LED,
}

const srcAddr = 0x01
