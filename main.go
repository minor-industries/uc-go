//go:build rp2040

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"
	"tinygo/exe/ir"
	"tinygo/leds"
	"tinygo/pio"
)

const blinkWrapTarget = 2
const blinkWrap = 7

var blinkProgram = pio.PIOProgram{
	Instructions: []uint16{
		0x80a0, //  0: pull   block
		0x6040, //  1: out    y, 32
		//     .wrap_target
		0xa022, //  2: mov    x, y
		0xe001, //  3: set    pins, 1
		0x0044, //  4: jmp    x--, 4
		0xa022, //  5: mov    x, y
		0xe000, //  6: set    pins, 0
		0x0047, //  7: jmp    x--, 7
		//     .wrap
	},
	Origin: -1,
}

func blinkProgramDefaultConfig(offset uint8) pio.PIOStateMachineConfig {
	cfg := pio.DefaultStateMachineConfig()
	cfg.SetWrap(offset+blinkWrapTarget, offset+blinkWrap)
	return cfg
}

// this is a raw helper function for use by the user which sets up the GPIO output, and configures the SM to output on a particular pin
func blinkProgramInit(sm *pio.PIOStateMachine, offset uint8, pin machine.Pin) {
	pin.Configure(machine.PinConfig{Mode: machine.PinPIO0})
	sm.SetConsecutivePinDirs(pin, 1, true)
	cfg := blinkProgramDefaultConfig(offset)
	cfg.SetSetPins(pin, 1)
	sm.Init(offset, &cfg)
}

func pioMain() {
	p := pio.PIO0
	p.Configure()

	offset := p.AddProgram(&blinkProgram)
	fmt.Printf("Loaded program at %d\n", offset)

	blinkPinForever(&p.StateMachines[0], offset, machine.LED, 3)
	blinkPinForever(&p.StateMachines[1], offset, machine.GPIO6, 4)
	blinkPinForever(&p.StateMachines[2], offset, machine.GPIO11, 1)
}

func main() {
	ir.Main()
	////pioMain()
	sm := leds.Setup()

	pixels := make([]color.RGBA, 150)
	pixels[0].R = 0x10
	pixels[1].G = 0x10
	pixels[2].B = 0x10

	for {
		fmt.Printf("tx\r\n")
		leds.Write(sm, pixels)
		time.Sleep(1000 * time.Millisecond)
	}
}

const clockHz = 133000000

func blinkPinForever(sm *pio.PIOStateMachine, offset uint8, pin machine.Pin, freq uint) {
	blinkProgramInit(sm, offset, pin)
	sm.SetEnabled(true)

	fmt.Printf("Blinking pin %d at %d Hz\n", pin, freq)
	sm.Tx(uint32(clockHz / (2 * freq)))
}
