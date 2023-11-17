package max31855

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"uc-go/pkg/spi"
)

type Thermocouple struct {
	spi *spi.SPI
	log func(string)
}

func NewThermocouple(
	spi *spi.SPI,
	log func(string),
) *Thermocouple {
	return &Thermocouple{
		spi: spi,
		log: log,
	}
}

func (tc *Thermocouple) Temperature() (float64, error) {
	var rx [4]byte

	err := tc.spi.TxSPI(
		[]byte{0, 0, 0, 0},
		rx[:],
	)
	if err != nil {
		return 0, errors.Wrap(err, "txspi")
	}

	tc.log("rx: " + hex.Dump(rx[:]))

	return -1, errors.New("not implemented")
}
