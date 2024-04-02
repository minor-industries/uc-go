package tempmon

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/minor-industries/uc-go/pkg/blikenlights"
	"github.com/minor-industries/uc-go/pkg/logger"
	rfm69_board "github.com/minor-industries/uc-go/pkg/rfm69-board"
	"github.com/minor-industries/uc-go/pkg/schema"
	"github.com/minor-industries/uc-go/pkg/spi"
	"github.com/minor-industries/uc-go/pkg/storage"
	"github.com/pkg/errors"
	"sync"
	"time"
	"tinygo.org/x/drivers/aht20"
)

const dstAddr = 2

var tStart time.Time

func Run(logs logger.Logger) error {
	tStart = time.Now()

	bl := blikenlights.NewLight(&blinker)
	go bl.Run()
	bl.Seq([]int{4, 4})

	<-time.After(2 * time.Second)
	fmt.Println("start")

	bl.Seq([]int{2, 2})

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
		bl,
		radio,
		sensor,
	)
	if err != nil {
		return errors.Wrap(err, "mainloop")
	}

	return errors.New("run exited")
}

func mainLoop(
	logs logger.Logger,
	bl *blikenlights.Light,
	radio *rfm69.Radio,
	sensor aht20.Device,
) error {
	radio.SetMode(rfm69.ModeSleep)

	readAndSend := func() error {
		bl.On()
		<-time.After(25 * time.Millisecond)
		bl.Off()

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

		afterStartupDelay := time.Now().Sub(tStart) > 20*time.Second

		if afterStartupDelay {
			pause()
		} else {
			<-time.After(5 * time.Second)
		}

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
	//arm.Asm("wfi")
	<-time.After(time.Second)
}

//// Wait enters WAIT (sleep) mode
//func Wait() {
//	// clear SLEEPDEEP bit to disable deep sleep
//	nxp.SystemControl.SCR.ClearBits(nxp.SystemControl_SCR_SLEEPDEEP)
//
//	// enter WAIT mode
//	arm.Asm("wfi")
//}
