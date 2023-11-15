package rfm69_board

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"sync/atomic"
	"time"
	"uc-go/pkg/spi"
)

type Board struct {
	spi  *spi.SPI
	rst  machine.Pin
	intr machine.Pin

	interruptCount uint32
	unhandledCount uint32

	interruptCh chan struct{}

	log func(s string)
}

func NewBoard(
	spi *spi.SPI,
	rst machine.Pin,
	intr machine.Pin,
	log func(s string),
) (*Board, error) {
	b := &Board{
		spi:         spi,
		rst:         rst,
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
	return b.spi.TxSPI(w, r)
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
