package main

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"time"
	"uc-go/pkg/blikenlights"
	rfm69_board "uc-go/pkg/rfm69-board"
	rfmCfg "uc-go/pkg/rfm69-board/cfg"
	"uc-go/pkg/spi"
)

type logger struct{}

func (l *logger) Log(s string) {
	fmt.Println(s)
}

func (l *logger) Error(err error) {
	fmt.Printf("error: %v\n", err)
}

func (l *logger) Rpc(s string, i interface{}) error {
	fmt.Println("rpc: " + s)
	return nil
}

type Cfg struct {
	led machine.Pin
}

func setupLeds(cfg *Cfg) *blikenlights.Light {
	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bl := blikenlights.NewLight(cfg.led)
	go bl.Run()
	bl.Seq([]int{2, 2})
	return bl
}

func run(log *logger) error {
	env := &Cfg{
		led: machine.PA07,
	}

	bl := setupLeds(env)

	<-time.After(5 * time.Second)
	bl.Seq([]int{4, 4})

	rfmSPILock := new(sync.Mutex)
	rfmSPI := spi.NewSPI(
		&spi.Config{
			Spi: &machine.SPI0,
			Config: &machine.SPIConfig{
				Frequency: 0,
				SCK:       machine.PA11,
				SDO:       machine.PA10,
				SDI:       machine.PA09,
				LSBFirst:  false,
				Mode:      0,
			},
			Cs: 0,
		},
		rfmSPILock,
	)

	_, err := rfm69_board.SetupRfm69(
		&rfmCfg.Config{
			NodeAddr: 100,
			TxPower:  20,
		},
		rfmSPI,
		&rfm69_board.PinCfg{
			Rst:  machine.PA15,
			Intr: machine.PA03,
		},
		func(s string) {
			log.Log(s)
		},
	)
	if err != nil {
		log.Error(errors.Wrap(err, "setup radio"))
	}

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		log.Log("hello")
	}

	return nil
}

func main() {
	log := &logger{}

	err := run(log)
	if err != nil {
		log.Error(errors.Wrap(err, "run exited with error"))
	} else {
		log.Log("run exited")
	}
}
