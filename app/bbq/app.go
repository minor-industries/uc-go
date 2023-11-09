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
)

const dstAddr = 2

type tcBoard struct {
	spiLock *sync.Mutex
	spi     *machine.SPI
	csn     machine.Pin
}

func (t *tcBoard) TxSPI(w, r []byte) error {
	t.spiLock.Lock()
	defer t.spiLock.Unlock()

	t.spi.Configure(machine.SPIConfig{
		Mode: 1,
		SCK:  machine.GPIO2,
		SDO:  machine.GPIO3,
		SDI:  machine.GPIO4,
	})

	t.csn.Low()
	err := t.spi.Tx(w, r)
	t.csn.High()

	return err
}

func Run(logs *rpc.Queue) error {
	stopLeds := make(chan struct{})
	go ledControl(stopLeds)
	go func() {
		<-time.After(5 * time.Second)
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

	cfg.Tc.Csn.Set(true)
	cfg.Tc.Csn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	tc := max31856.NewMAX31856(&tcBoard{
		spi:     cfg.Tc.Spi,
		spiLock: spiLock,
		csn:     cfg.Tc.Csn,
	}, log)

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

	<-time.After(time.Second)

	tc.Init()
	logs.Log("tc init complete")

	err = mainLoop(
		logs,
		radio,
		rand.New(rand.NewSource(int64(envSnapshot.NodeAddr))),
		tc,
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
	tc *max31856.MAX31856,
) error {
	readAndSend := func() error {
		return errors.New("read and send not implemented")
	}

	for {
		if err := readAndSend(); err != nil {
			t := tc.Temperature()
			logs.Log(fmt.Sprintf("temp = %f", t))
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}
