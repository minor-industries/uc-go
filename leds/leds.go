package leds

import (
	"fmt"
	"machine"
	"time"
	"tinygo/pio"
)

func tightLoopContents() {}

func Run() {
	p := pio.PIO0
	p.Configure()

	offset := p.AddProgram(&ws2812Program)
	fmt.Printf("Loaded program at %d\n", offset)

	sm := &p.StateMachines[0]
	ws2812ProgramInit(sm, offset, machine.GP0)
	sm.SetEnabled(true)

	const smTxFullMask = 0x1

	for {
		fmt.Printf("tx\r\n")
		for i := 0; i < 150; i++ {
			for p.Device.GetFSTAT_TXFULL()&smTxFullMask != 0 {
				tightLoopContents()
			}
			sm.Tx(0x00102000)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
