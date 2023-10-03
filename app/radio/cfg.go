package radio

import "machine"

type pinCfg struct {
	// rfm
	spi *machine.SPI

	rst  machine.Pin
	intr machine.Pin

	sck machine.Pin
	sdo machine.Pin
	sdi machine.Pin
	csn machine.Pin

	// i2c
	i2c *machine.I2C
	sda machine.Pin
	scl machine.Pin

	// misc
	led machine.Pin
}
