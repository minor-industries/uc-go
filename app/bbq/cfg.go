package bbq

import (
	"machine"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

type ThermocoupleCfg struct {
	Name string
	Spi  *spi.Config
}

type BoardCfg struct {
	// i2c
	Rfm rfm69_board.PinCfg

	// misc
	led machine.Pin

	Tcs []*ThermocoupleCfg
}
