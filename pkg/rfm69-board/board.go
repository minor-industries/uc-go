package rfm69_board

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"sync/atomic"
	"time"
)

type Board struct {
	spi     *machine.SPI
	spiLock *sync.Mutex
	spiCfg  *machine.SPIConfig

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
	spiCfg *machine.SPIConfig,
	spiLock *sync.Mutex,
	rst machine.Pin,
	csn machine.Pin,
	intr machine.Pin,
	log func(s string),
) (*Board, error) {
	b := &Board{
		spi:         spi,
		spiCfg:      spiCfg,
		spiLock:     spiLock,
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
	b.spiLock.Lock()
	defer b.spiLock.Unlock()

	err := b.spi.Configure(*b.spiCfg)
	if err != nil {
		return errors.Wrap(err, "configure spi")
	}

	b.csn.Low()
	err = b.spi.Tx(w, r)
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
