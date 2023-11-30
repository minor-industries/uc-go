package tempmon

import (
	"device/arm"
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"time"
	"tinygo.org/x/drivers/aht20"
	"uc-go/pkg/blikenlights"
	"uc-go/pkg/logger"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/schema"
	"uc-go/pkg/spi"
	"uc-go/pkg/storage"
)

const dstAddr = 2

func Run(logs logger.Logger) error {
	bl := blikenlights.NewLight(&blinker)
	go bl.Run()
	bl.Seq([]int{2, 2})

	<-time.After(10 * time.Second)
	fmt.Println("start")

	bl.Seq([]int{2, 4})

	lfs, err := storage.Setup(logs)
	if err != nil {
		return errors.Wrap(err, "setup storage")
	}

	env, err := rfm69_board.LoadConfig(logs, lfs, false)
	if err != nil {
		return errors.Wrap(err, "load config")
	}
	envSnapshot := env.SnapShot()
	envSnapshot.NodeAddr = 0xD0 // TODO: need to fix this config stuff

	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	stopLeds := make(chan struct{})
	go ledControl(stopLeds)
	go func() {
		<-time.After(5 * time.Second)
		close(stopLeds)
	}()

	fmt.Printf("address is 0x%02x\n", envSnapshot.NodeAddr)

	i2c := cfg.i2c

	err = i2c.Configure(*cfg.i2cCfg)
	if err != nil {
		return errors.Wrap(err, "configure i2c")
	}

	sensor := aht20.New(i2c)
	sensor.Configure()
	sensor.Reset()

	log := func(s string) {
		logs.Log(s)
	}

	spiLock := new(sync.Mutex)
	radioSpi := spi.NewSPI(cfg.Rfm.Spi, spiLock)

	radio, err := rfm69_board.SetupRfm69(
		&envSnapshot,
		&cfg.Rfm,
		radioSpi,
		log,
	)
	if err != nil {
		return errors.Wrap(err, "setup radio")
	}

	err = mainLoop(
		logs,
		radio,
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
	logs logger.Logger,
	radio *rfm69.Radio,
	sensor aht20.Device,
) error {
	radio.SetMode(rfm69.ModeSleep)

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

		if err := rfm69_board.SendMsg(radio, dstAddr, 1, body); err != nil {
			return errors.Wrap(err, "send msg")
		}

		radio.SetMode(rfm69.ModeSleep)

		return nil
	}

	for {
		if err := readAndSend(); err != nil {
			logs.Error(err)
		}

		pause()
	}
}

func pause() {
	end := time.Now().Add(5 * time.Second)
	for {
		SleepDeep()
		if time.Now().After(end) {
			return
		}
	}
}

// SleepDeep enters STOP (deep sleep) mode
func SleepDeep() {
	// set SLEEPDEEP to enable deep sleep

	//arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	arm.Asm("wfi")
}

//// Wait enters WAIT (sleep) mode
//func Wait() {
//	// clear SLEEPDEEP bit to disable deep sleep
//	nxp.SystemControl.SCR.ClearBits(nxp.SystemControl_SCR_SLEEPDEEP)
//
//	// enter WAIT mode
//	arm.Asm("wfi")
//}
