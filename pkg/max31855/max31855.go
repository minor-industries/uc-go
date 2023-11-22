package max31855

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
	for _, b := range rx {
		fmt.Printf("%08b ", b)
	}
	fmt.Println("")

	if rx[3]&0x01 != 0 {
		return 0, errors.New("thermocouple not connected")
	}

	if rx[3]&0x02 != 0 {
		return 0, errors.New("short circuit to ground")
	}

	if rx[3]&0x04 != 0 {
		return 0, errors.New("short circuit to power")
	}

	if rx[1]&0x01 != 0 {
		return 0, errors.New("fault")
	}

	// TODO: handle/test negative temperatures
	hotTemp := float64(binary.BigEndian.Uint16(rx[0:2])>>2) / 4.0
	refTemp := float64(binary.BigEndian.Uint16(rx[2:4])>>4) / 16.0

	fmt.Println(hotTemp, refTemp)

	return hotTemp, nil
}
