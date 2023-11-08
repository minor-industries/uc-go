package radio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"math/rand"
	"time"
	"tinygo.org/x/drivers/aht20"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
	cfg3 "uc-go/pkg/rfm69-board/cfg"
	"uc-go/pkg/schema"
)

func Run(logs *rpc.Queue) error {
	env, err := rfm69_board.LoadConfig(logs)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

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

	envSnapshot := env.SnapShot()

	err = runRadio(
		logs,
		&envSnapshot,
		log,
		sensor,
	)
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
	env *cfg3.Config,
	log func(s string),
	sensor aht20.Device,
) error {
	randSource := rand.New(rand.NewSource(int64(env.NodeAddr)))

	radio, err := rfm69_board.SetupRfm69(env, &cfg.Rfm, log)
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
			env.NodeAddr,
			env.TxPower,
			bodyBuf.Bytes(),
		); err != nil {
			logs.Error(errors.Wrap(err, "send frame"))
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}
