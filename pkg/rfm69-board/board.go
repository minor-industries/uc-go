package rfm69_board

import "machine"

type Board struct {
	spi *machine.SPI
	rst machine.Pin
	csn machine.Pin
}

func NewBoard(spi *machine.SPI, rst machine.Pin, csn machine.Pin) *Board {
	return &Board{spi: spi, rst: rst, csn: csn}
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
