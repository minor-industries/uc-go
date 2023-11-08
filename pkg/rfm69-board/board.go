package rfm69_board

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"sync/atomic"
	"time"
	"uc-go/pkg/protocol/rpc"
	cfg3 "uc-go/pkg/rfm69-board/cfg"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
)

type Board struct {
	spi  *machine.SPI
	rst  machine.Pin
	csn  machine.Pin
	intr machine.Pin

	interruptCount uint32
	unhandledCount uint32

	interruptCh chan struct{}

	log func(s string)
}

func NewBoard(
	spi *machine.SPI,
	rst machine.Pin,
	csn machine.Pin,
	intr machine.Pin,
	log func(s string),
) (*Board, error) {
	b := &Board{
		spi:         spi,
		rst:         rst,
		csn:         csn,
		intr:        intr,
		interruptCh: make(chan struct{}),
		log:         log,
	}

	b.intr.Configure(machine.PinConfig{Mode: machine.PinInput})
	if err := b.intr.SetInterrupt(machine.PinRising, b.handleInterrupt); err != nil {
		return nil, errors.Wrap(err, "set interrupt")
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			log(fmt.Sprintf(
				"interrupt count = %d, unhandled count = %d",
				atomic.LoadUint32(&b.interruptCount),
				atomic.LoadUint32(&b.unhandledCount),
			))
		}
	}()

	return b, nil
}

func (b *Board) TxSPI(w, r []byte) error {
	b.csn.Low()
	err := b.spi.Tx(w, r)
	b.csn.High()
	return err
}

func (b *Board) Reset(b2 bool) error {
	b.rst.Set(b2)
	return nil
}

func (b *Board) WaitForD0Edge() {
	<-b.interruptCh
}

func (b *Board) handleInterrupt(pin machine.Pin) {
	atomic.AddUint32(&b.interruptCount, 1)
	select {
	case b.interruptCh <- struct{}{}:
	default:
		atomic.AddUint32(&b.unhandledCount, 1)
	}
}

type PinCfg struct {
	// rfm
	Spi *machine.SPI

	Rst  machine.Pin
	Intr machine.Pin

	Sck machine.Pin
	Sdo machine.Pin
	Sdi machine.Pin
	Csn machine.Pin
}

func SetupRfm69(
	cfg *PinCfg,
	log func(s string),
) (*rfm69.Radio, error) {
	rst := cfg.Rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := cfg.Spi
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode0,
		SCK:  cfg.Sck,
		SDO:  cfg.Sdo,
		SDI:  cfg.Sdi,
	})

	log("hello")

	//return nil, err

	if err != nil {
		return nil, errors.Wrap(err, "configure SPI")
	} else {
		log("setup SPI")
	}

	CSn := cfg.Csn

	CSn.Set(true)
	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CSn.Set(true)

	board, err := NewBoard(
		spi,
		rst,
		CSn,
		cfg.Intr,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new board")
	}

	radio := rfm69.NewRadio(board, log)

	if err := radio.Setup(rfm69.RF69_433MHZ); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	return radio, nil
}

const (
	configFile = "/radio-cfg.msgp"

	initialTxPower  = 20
	initialNodeAddr = 0xee
)

func LoadConfig(logs *rpc.Queue) (
	*util.SyncConfig[cfg3.Config],
	error,
) {
	lfs, err := storage.Setup(logs)
	if err != nil {
		return nil, errors.Wrap(err, "setup storage")
	}

	if lfs == nil {
		return nil, errors.New("no lfs")
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
		return nil, errors.Wrap(err, "load config")
	}

	env := util.NewSyncConfig(*config)
	return env, nil
}
