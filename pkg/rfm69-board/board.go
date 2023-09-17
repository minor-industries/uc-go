package rfm69_board

import (
	"github.com/pkg/errors"
	"machine"
	"sync/atomic"
)

type Board struct {
	spi  *machine.SPI
	rst  machine.Pin
	csn  machine.Pin
	intr machine.Pin

	interruptCount uint32
}

func NewBoard(
	spi *machine.SPI,
	rst machine.Pin,
	csn machine.Pin,
	intr machine.Pin,
) (*Board, error) {
	b := &Board{
		spi:  spi,
		rst:  rst,
		csn:  csn,
		intr: intr,
	}

	b.intr.Configure(machine.PinConfig{Mode: machine.PinInput})
	if err := b.intr.SetInterrupt(machine.PinRising, b.handleInterrupt); err != nil {
		return nil, errors.Wrap(err, "set interrupt")
	}

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
	select {} // TODO
}

func (b *Board) handleInterrupt(pin machine.Pin) {
	atomic.AddUint32(&b.interruptCount, 1)
}
