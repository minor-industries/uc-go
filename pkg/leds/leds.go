//go:build rp2040

package leds

import (
	"fmt"
	"image/color"
	"machine"
	"sync/atomic"
	"uc-go/pkg/pio"
)

func tightLoopContents() {}

var TxFullCounter int64

func Setup() *pio.PIOStateMachine {
	p := pio.PIO0
	p.Configure()

	offset := p.AddProgram(&ws2812Program)
	fmt.Printf("Loaded program at %d\n", offset)

	sm := &p.StateMachines[0]
	ws2812ProgramInit(sm, offset, machine.GP0)
	sm.SetEnabled(true)

	return sm
}

func Write(sm *pio.PIOStateMachine, pixels []color.RGBA) {
	const smTxFullMask = 0x1

	for _, pixel := range pixels {
		for sm.PIO.Device.GetFSTAT_TXFULL()&smTxFullMask != 0 {
			atomic.AddInt64(&TxFullCounter, 1)
			tightLoopContents()
		}
		r := uint32(pixel.R)
		g := uint32(pixel.G)
		b := uint32(pixel.B)

		v := g<<24 + r<<16 + b<<8
		sm.Tx(v)
	}
}
