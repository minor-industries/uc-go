package leds

import (
	"fmt"
	"machine"
	"time"
	"tinygo/pio"
)

func tightLoopContents() {}

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

func Write(sm *pio.PIOStateMachine) {
	const smTxFullMask = 0x1

	for {
		fmt.Printf("tx\r\n")
		for i := 0; i < 150; i++ {
			for sm.PIO.Device.GetFSTAT_TXFULL()&smTxFullMask != 0 {
				tightLoopContents()
			}
			sm.Tx(0x00101100)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
