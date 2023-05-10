package leds

import (
	"fmt"
	"machine"
	"time"
	"tinygo/pio"
)

func Run() {
	p := pio.PIO0
	p.Configure()

	offset := p.AddProgram(&ws2812Program)
	fmt.Printf("Loaded program at %d\n", offset)

	sm := &p.StateMachines[0]
	ws2812ProgramInit(sm, offset, machine.GP0)
	sm.SetEnabled(true)

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("tx\r\n")
		sm.Tx(0x00404000)
		sm.Tx(0x00404000)
		sm.Tx(0x00404000)
		sm.Tx(0x00404000)
	}
}
