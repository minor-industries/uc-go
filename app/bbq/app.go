package bbq

import (
	"encoding/hex"
	"fmt"
	"github.com/minor-industries/max31856"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"time"
	"uc-go/pkg/blikenlights"
	"uc-go/pkg/logger"
	"uc-go/pkg/max31855"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/schema"
	"uc-go/pkg/spi"
)

const dstAddr = 2

func Run(logs logger.Logger) error {
	bl := blikenlights.NewLight(cfg.led)
	go bl.Run()
	bl.Seq([]int{2, 2})

	//go func() {
	//	<-time.After(5 * time.Minute)
	//	close(stopLeds)
	//}()

	time.Sleep(2 * time.Second)
	fmt.Println("starting")

	env, err := rfm69_board.LoadConfig(logs)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	envSnapshot := env.SnapShot()
	_ = envSnapshot // TODO

	spiLock := new(sync.Mutex)

	// need to build all SPIs before using them to set CS lines
	rfmSpi := spi.NewSPI(cfg.Rfm.Spi, spiLock)

	log := func(s string) {
		logs.Log(s)
	}

	tcNames := make([]string, len(cfg.Tcs))
	tcs := map[string]*max31856.MAX31856{}

	for i, tc := range cfg.Tcs {
		tcNames[i] = tc.Name
		tcs[tc.Name] = max31856.NewMAX31856(
			spi.NewSPI(tc.Spi, spiLock),
			log,
		)
	}

	for _, name := range tcNames {
		tc := tcs[name]
		err := tc.Init()
		if err != nil {
			logs.Error(errors.Wrap(err, fmt.Sprintf("tc [%s] init error", name)))
		} else {
			logs.Log(fmt.Sprintf("tc [%s] init", name))
		}
	}

	{
		bl.Seq([]int{4, 4})
		once := sync.Once{}

		for _, name := range tcNames {
			tc := tcs[name]
			t := tc.Temperature()
			if t != 0 {
				once.Do(func() {
					bl.Seq([]int{32 - 12, 4, 4, 4})
				})
			}
			logs.Log(fmt.Sprintf("ABC %s: tc [%s] temp = %.02f", time.Now().String(), name, t))
		}
	}

	i2c := machine.I2C0
	err = i2c.Configure(machine.I2CConfig{})
	if err != nil {
		logs.Error(errors.Wrap(err, "configure i2c"))
	} else {
		logs.Log("configured i2c")
	}

	{
		var rx [2]byte
		err := i2c.ReadRegister(0x60, 0x00, rx[:])
		if err != nil {
			logs.Error(errors.Wrap(err, "read register"))
		} else {
			logs.Log("read: " + hex.Dump(rx[:]))
		}
	}

	spiX := spi.NewSPI(&spi.Config{
		Spi: &machine.SPI0,
		Config: &machine.SPIConfig{
			SCK:  machine.SPI0_SCK_PIN,
			SDO:  machine.SPI0_SDO_PIN,
			SDI:  machine.SPI0_SDI_PIN,
			Mode: 1,
		},
		Cs: machine.A2,
	}, spiLock)

	tc2 := max31855.NewThermocouple(spiX, log)

	radio, err := rfm69_board.SetupRfm69(
		&envSnapshot,
		&cfg.Rfm,
		rfmSpi,
		log,
	)
	if err != nil {
		logs.Error(err)
		//return errors.Wrap(err, "rfm69")
	}

	err = mainLoop(
		logs,
		radio,
		tcNames,
		tcs,
		tc2,
	)
	if err != nil {
		return errors.Wrap(err, "mainloop")
	}

	return errors.New("run exited")
}

func mainLoop(
	logs logger.Logger,
	radio *rfm69.Radio,
	tcNames []string,
	tcs map[string]*max31856.MAX31856,
	tc2 *max31855.Thermocouple,
) error {

	readAndSend := func() error {
		for _, name := range tcNames {
			tc := tcs[name]
			t := tc.Temperature()
			logs.Log(fmt.Sprintf("tc [%s] temp = %.02f", name, t))

			tc2.Temperature()

			if radio == nil {
				return errors.New("no radio")
			}

			desc := [16]byte{}
			copy(desc[:], name)

			if err := rfm69_board.SendMsg(radio, dstAddr, 2, &schema.ThermocoupleData{
				Temperature: float32(t),
				Description: desc,
			}); err != nil {
				return errors.Wrap(err, "send msg")
			}
		}
		return nil
	}

	for {
		if err := readAndSend(); err != nil {
			logs.Error(errors.Wrap(err, "read and send"))
		}

		// TODO: need random sleep here (with low memory usage)
		sleep := 5000 * time.Millisecond
		time.Sleep(sleep)
	}
}
