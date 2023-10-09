package radio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"math/rand"
	"time"
	"tinygo.org/x/drivers/aht20"
	cfg3 "uc-go/app/radio/cfg"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/schema"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
)

const (
	configFile = "/radio-cfg.msgp"

	initialTxPower  = 20
	initialNodeAddr = 0xee
)

func Run(logs *rpc.Queue) error {
	lfs, err := storage.Setup(logs)
	if err != nil {
		return errors.Wrap(err, "setup storage")
	}

	if lfs == nil {
		return errors.New("no lfs")
	}

	config, err := storage.LoadConfig[*cfg3.Config](
		lfs,
		logs,
		configFile,
		&cfg3.Config{
			NodeAddr: initialNodeAddr,
			TxPower:  initialTxPower,
		},
	)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	env := util.NewSyncConfig(*config)

	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	stopLeds := make(chan struct{})
	go ledControl(stopLeds)
	go func() {
		<-time.After(5 * time.Second)
		close(stopLeds)
	}()

	i2c := cfg.i2c

	err = i2c.Configure(machine.I2CConfig{
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

	err = runRadio(logs, env, log, sensor)
	if err != nil {
		return errors.Wrap(err, "run radio")
	}

	return errors.New("run exited")
}

func ledControl(done <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	val := true

	for {
		select {
		case <-ticker.C:
			cfg.led.Set(val)
			val = !val
		case <-done:
			return
		}
	}
}

func runRadio(
	logs *rpc.Queue,
	env *util.SyncConfig[cfg3.Config],
	log func(s string),
	sensor aht20.Device,
) error {
	envSnapshot := env.SnapShot()
	randSource := rand.New(rand.NewSource(int64(envSnapshot.NodeAddr)))

	radio, err := setupRfm69(log)
	if err != nil {
		return errors.Wrap(err, "rfm69")
	}

	for {
		err := sensor.Read()
		if err != nil {
			logs.Error(errors.Wrap(err, "read sensor"))
			continue
		}
		t := sensor.Celsius()
		logs.Log(fmt.Sprintf("temperature = %0.01fC %0.01fF", t, (t*9/5)+32))

		body := &schema.SensorData{
			Temperature:      t,
			RelativeHumidity: sensor.RelHumidity(),
		}

		bodyBuf := bytes.NewBuffer(nil)
		bodyBuf.WriteByte(1) // message ID
		if err := binary.Write(bodyBuf, binary.LittleEndian, body); err != nil {
			return errors.Wrap(err, "encode body")
		}

		if err := radio.SendFrame(
			2,
			envSnapshot.NodeAddr,
			envSnapshot.TxPower,
			bodyBuf.Bytes(),
		); err != nil {
			logs.Error(errors.Wrap(err, "send frame"))
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}

func setupRfm69(log func(s string)) (*rfm69.Radio, error) {
	rst := cfg.rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := cfg.spi
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode0,
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

	board, err := rfm69_board.NewBoard(
		spi,
		rst,
		CSn,
		cfg.intr,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new board")
	}

	//return nil, errors.New("error-01")

	radio := rfm69.NewRadio(board, log)

	if err := radio.Setup(rfm69.RF69_915MHZ); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	return radio, nil
}
