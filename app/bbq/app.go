package bbq

import (
	"fmt"
	"github.com/minor-industries/max31856"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"math/rand"
	"sync"
	"time"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
	"uc-go/pkg/spi"
)

const dstAddr = 2

func Run(logs *rpc.Queue) error {
	stopLeds := make(chan struct{})
	go ledControl(stopLeds)
	go func() {
		<-time.After(5 * time.Minute)
		close(stopLeds)
	}()

	env, err := rfm69_board.LoadConfig(logs)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	envSnapshot := env.SnapShot()
	spiLock := new(sync.Mutex)

	log := func(s string) {
		logs.Log(s)
	}

	tcNames := []string{
		cfg.Tc0.Name,
		cfg.Tc1.Name,
	}

	tcs := map[string]*max31856.MAX31856{
		cfg.Tc0.Name: max31856.NewMAX31856(
			spi.NewSPI(cfg.Tc0.Spi, spiLock),
			log,
		),
		cfg.Tc1.Name: max31856.NewMAX31856(
			spi.NewSPI(cfg.Tc0.Spi, spiLock),
			log,
		),
	}

	for _, name := range tcNames {
		tc := tcs[name]
		err := tc.Init()
		if err != nil {
			logs.Error(errors.Wrap(err, fmt.Sprintf("tc [%s] init error", name)))
		} else {
			logs.Log(fmt.Sprintf("tc [%s] init"))
		}
	}

	radio, err := rfm69_board.SetupRfm69(
		&envSnapshot,
		&cfg.Rfm,
		spiLock,
		log,
	)
	if err != nil {
		logs.Error(err)
		//return errors.Wrap(err, "rfm69")
	}

	err = mainLoop(
		logs,
		radio,
		rand.New(rand.NewSource(int64(envSnapshot.NodeAddr))),
		tcNames,
		tcs,
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
	tcNames []string,
	tcs map[string]*max31856.MAX31856,
) error {
	readAndSend := func() error {
		for _, name := range tcNames {
			tc := tcs[name]
			t := tc.Temperature()
			logs.Log(fmt.Sprintf("tc [%s] temp = %.02f", name, t))
		}
		return nil
	}

	for {
		if err := readAndSend(); err != nil {
			logs.Error(errors.Wrap(err, "read and send"))
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}
