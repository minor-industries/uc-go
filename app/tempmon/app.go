package tempmon

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"math/rand"
	"sync"
	"time"
	"tinygo.org/x/drivers/aht20"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/schema"
)

const dstAddr = 2

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
	spiLock := new(sync.Mutex)

	radio, err := rfm69_board.SetupRfm69(
		&envSnapshot,
		&cfg.Rfm,
		spiLock,
		log,
	)
	if err != nil {
		return errors.Wrap(err, "setup radio")
	}

	err = mainLoop(
		logs,
		radio,
		rand.New(rand.NewSource(int64(envSnapshot.NodeAddr))),
		sensor,
	)
	if err != nil {
		return errors.Wrap(err, "mainloop")
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

func mainLoop(
	logs *rpc.Queue,
	radio *rfm69.Radio,
	randSource *rand.Rand,
	sensor aht20.Device,
) error {
	readAndSend := func() error {
		err := sensor.Read()
		if err != nil {
			return errors.Wrap(err, "read sensor")
		}
		t := sensor.Celsius()
		logs.Log(fmt.Sprintf("temperature = %0.01fC %0.01fF", t, (t*9/5)+32))

		body := &schema.SensorData{
			Temperature:      t,
			RelativeHumidity: sensor.RelHumidity(),
		}

		return rfm69_board.SendMsg(radio, dstAddr, 1, body)
	}

	for {
		if err := readAndSend(); err != nil {
			logs.Error(err)
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}