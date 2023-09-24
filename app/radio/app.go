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

func Run(logs *rpc.Queue) error {
	i2c := machine.I2C0

	err := i2c.Configure(machine.I2CConfig{
		SDA: machine.GP0,
		SCL: machine.GP1,
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

	radio, err := setupRfm69(log)
	if err != nil {
		return errors.Wrap(err, "rfm69")
	}

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

	return errors.New("run exited")
}

func setupRfm69(log func(s string)) (*rfm69.Radio, error) {
	rst := machine.GP6
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := machine.SPI0
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode3,
		SCK:  machine.GP2,
		SDO:  machine.GP3,
		SDI:  machine.GP4,
	})
	if err != nil {
		return nil, errors.Wrap(err, "configure SPI")
	} else {
		log("setup SPI")
	}

	CSn := machine.GP5
	CSn.Set(true)
	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CSn.Set(true)

	board, err := rfm69_board.NewBoard(
		spi,
		rst,
		CSn,
		machine.GP7,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new board")
	}

	radio := rfm69.NewRadio(board, log)

	if err := radio.Setup(rfm69.RF69_915MHZ); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	return radio, nil
}
