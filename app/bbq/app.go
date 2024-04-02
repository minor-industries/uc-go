package bbq

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/minor-industries/uc-go/pkg/blikenlights"
	"github.com/minor-industries/uc-go/pkg/logger"
	"github.com/minor-industries/uc-go/pkg/max31855"
	"github.com/minor-industries/uc-go/pkg/mcp9600"
	rfm69_board "github.com/minor-industries/uc-go/pkg/rfm69-board"
	"github.com/minor-industries/uc-go/pkg/schema"
	"github.com/minor-industries/uc-go/pkg/spi"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"time"
)

const dstAddr = 2

type Thermometer interface {
	Temperature() (float64, error)
}

type App struct {
	i2c    *machine.I2C
	rfmSpi *spi.SPI
	tcSpi  *spi.SPI
	logs   logger.Logger
	radio  *rfm69.Radio
	tcs    map[string]Thermometer
}

func (app *App) setupI2C() {
	app.i2c = machine.I2C0
	err := app.i2c.Configure(machine.I2CConfig{})
	if err != nil {
		app.logs.Error(errors.Wrap(err, "configure i2c"))
	} else {
		app.logs.Log("configured i2c")
	}
}

func Run(logs logger.Logger) error {
	app := &App{
		logs: logs,
	}

	app.setupLeds()

	time.Sleep(2 * time.Second)
	fmt.Println("starting")

	app.setupI2C()
	app.setupSPI()

	app.setupTCs()

	err := app.setupRadio()
	if err != nil {
		app.logs.Error(errors.Wrap(err, "setup radio"))
	}

	err = app.mainLoop()
	if err != nil {
		return errors.Wrap(err, "mainloop")
	}

	return errors.New("run exited")
}

func (app *App) setupRadio() error {
	env, err := rfm69_board.LoadConfig(app.logs)
	if err != nil {
		return errors.Wrap(err, "load radio config")
	}
	envSnapshot := env.SnapShot()
	app.radio, err = rfm69_board.SetupRfm69(
		&envSnapshot,
		&cfg.Rfm,
		app.rfmSpi,
		app.log,
	)
	if err != nil {
		return errors.Wrap(err, "setup radio")
	}
	return nil
}

func (app *App) log(s string) {
	app.logs.Log(s)
}

func (app *App) setupLeds() {
	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bl := blikenlights.NewLight(cfg.led)
	go bl.Run()
	bl.Seq([]int{2, 2})

	//go func() {
	//	<-time.After(5 * time.Minute)
	//	close(stopLeds)
	//}()
}

func (app *App) setupSPI() {
	spiLock := new(sync.Mutex)

	// need to build all SPIs before using them to set CS lines
	app.rfmSpi = spi.NewSPI(cfg.Rfm.Spi, spiLock)
	app.tcSpi = spi.NewSPI(&spi.Config{
		Spi: &machine.SPI0,
		Config: &machine.SPIConfig{
			SCK:  machine.SPI0_SCK_PIN,
			SDO:  machine.SPI0_SDO_PIN,
			SDI:  machine.SPI0_SDI_PIN,
			Mode: 1,
		},
		Cs: machine.A0,
	}, spiLock)
}

func (app *App) mainLoop() error {
	readAndSend := func() error {
		for name, tc := range app.tcs {
			t, err := tc.Temperature()
			if err != nil {
				app.logs.Error(errors.Wrap(err, "reading temperature from "+name))
				continue
			}
			app.logs.Log(fmt.Sprintf("tc [%s] temp = %.02f", name, t))

			if app.radio == nil {
				return errors.New("no radio")
			}

			desc := [16]byte{}
			copy(desc[:], name)

			if err := rfm69_board.SendMsg(
				app.radio,
				dstAddr,
				2,
				&schema.ThermocoupleData{
					Temperature: float32(t),
					Description: desc,
				},
			); err != nil {
				return errors.Wrap(err, "send msg")
			}
		}
		return nil
	}

	for {
		if err := readAndSend(); err != nil {
			app.logs.Error(errors.Wrap(err, "read and send"))
		}

		// TODO: need random sleep here (with low memory usage)
		sleep := 5000 * time.Millisecond
		time.Sleep(sleep)
	}
}

func (app *App) setupTCs() {
	app.tcs = map[string]Thermometer{
		"max": max31855.NewThermocouple(app.tcSpi, app.log),
		"mcp": mcp9600.NewThermocouple(app.log, app.i2c, 0x60),
	}
}
