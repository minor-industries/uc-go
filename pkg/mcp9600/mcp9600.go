package mcp9600

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"machine"
)

type Thermocouple struct {
	log  func(string)
	i2c  *machine.I2C
	addr byte
}

func NewThermocouple(
	log func(string),
	i2c *machine.I2C,
	addr byte,
) *Thermocouple {
	return &Thermocouple{
		log:  log,
		i2c:  i2c,
		addr: addr,
	}
}

func (t *Thermocouple) Temperature() (float64, error) {
	var rx [2]byte
	err := t.i2c.ReadRegister(t.addr, 0x00, rx[:])
	if err != nil {
		return 0, errors.Wrap(err, "read register")
	} else {
		t.log("read: " + hex.Dump(rx[:]))
		t := float64(binary.BigEndian.Uint16(rx[:])) * 0.0625
		if rx[0]&0x80 != 0 {
			t -= 4096
		}
		fmt.Println("mpc t", t)
		return t, nil
	}
}
