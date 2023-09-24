package radio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"time"
	"tinygo.org/x/drivers/aht20"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
)

type SensorData struct {
	Temperature      float32 // celsius
	RelativeHumidity float32
}

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
}

//var pico = pinCfg{
//	spi: machine.SPI0,
//
//	rst:  machine.GP6,
//	intr: machine.GP7,
//
//	sck: machine.GP2,
//	sdo: machine.GP3,
//	sdi: machine.GP4,
//	csn: machine.GP5,
//
//	i2c: machine.I2C0,
//	sda: machine.GP0,
//	scl: machine.GP1,
//}

var featherRp2040 = pinCfg{
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
}

var cfg = featherRp2040

func Run(logs *rpc.Queue) error {
	i2c := cfg.i2c

	err := i2c.Configure(machine.I2CConfig{
		SDA: cfg.sda,
		SCL: cfg.scl,
	})
	if err != nil {
		return errors.Wrap(err, "configure i2c")
	}

	sensor := aht20.New(i2c)
	sensor.Configure()
	sensor.Reset()

	log := func(s string) {
		logs.Log(s)
	}

	err = runRadio(logs, log, sensor)
	if err != nil {
		return errors.Wrap(err, "run radio")
	}

	return errors.New("run exited")
}

func runRadio(
	logs *rpc.Queue,
	log func(s string),
	sensor aht20.Device,
) error {
	radio, err := setupRfm69(log)
	if err != nil {
		return errors.Wrap(err, "rfm69")
	}

	return errors.New("there-01")

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err := sensor.Read()
		if err != nil {
			logs.Error(errors.Wrap(err, "read sensor"))
			continue
		}
		t := sensor.Celsius()
		logs.Log(fmt.Sprintf("temperature = %0.01fC %0.01fF", t, (t*9/5)+32))

		body := &SensorData{
			Temperature:      t,
			RelativeHumidity: sensor.RelHumidity(),
		}

		bodyBuf := bytes.NewBuffer(nil)
		bodyBuf.WriteByte(1) // message ID
		if err := binary.Write(bodyBuf, binary.LittleEndian, body); err != nil {
			return errors.Wrap(err, "encode body")
		}

		if err := radio.SendFrame(2, 1, bodyBuf.Bytes()); err != nil {
			logs.Error(errors.Wrap(err, "send frame"))
		}
	}

	return nil
}

func setupRfm69(log func(s string)) (*rfm69.Radio, error) {
	rst := cfg.rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := cfg.spi
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode3,
		SCK:  cfg.sck,
		SDO:  cfg.sdo,
		SDI:  cfg.sdi,
	})

	log("hello")

	//return nil, err

	if err != nil {
		return nil, errors.Wrap(err, "configure SPI")
	} else {
		log("setup SPI")
	}

	CSn := cfg.csn

	CSn.Set(true)

	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//machine.GPIO16.Configure(machine.PinConfig{Mode: machine.PinOutput})

	CSn.Set(true)

	_, err = rfm69_board.NewBoard(
		spi,
		rst,
		CSn,
		cfg.intr,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new board")
	}

	return nil, errors.New("error-01")
	//
	//radio := rfm69.NewRadio(board, log)
	//
	//if err := radio.Setup(rfm69.RF69_915MHZ); err != nil {
	//	return nil, errors.Wrap(err, "setup")
	//}
	//
	//return radio, nil
}
