//go:build rp2040

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"
	"tinygo/bounce"
	"tinygo/cfg"
	"tinygo/exe/ir"
	"tinygo/leds"
	"tinygo/pio"
	"tinygo/strip"
)

const blinkWrapTarget = 2
const blinkWrap = 7

const (
	ledMaxLevel = 0.5 // brightness level of NeoPxels (0~1)
)

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

	runLeds(sm)
}

func runLeds(sm *pio.PIOStateMachine) {
	pixels := make([]color.RGBA, 150)
	pixels[0].R = 0x10
	pixels[1].G = 0x10
	pixels[2].B = 0x10

	strip := strip.NewStrip(&cfg.Cfg{
		NumLeds:    150,
		StartIndex: 0,
		Length:     5.0,
	})
	sim := bounce.Bounce(&bounce.App{Strip: strip})

	tick := 30 * time.Millisecond

	for range time.NewTicker(tick).C {
		sim.Tick(
			0,
			tick.Seconds(),
		)

		writeColors(sm, pixels, strip)
	}
}

func clamp(min, x, max float64) float64 {
	if x < min {
		return min
	}

	if x > max {
		return max
	}

	return x
}

func writeColors(
	sm *pio.PIOStateMachine,
	pixels []color.RGBA,
	st *strip.Strip,
) {
	st.Each(func(i int, led *strip.Led) {
		pixels[i].R = uint8(clamp(0, led.R, 1.0) * ledMaxLevel * 255.0)
		pixels[i].G = uint8(clamp(0, led.G, 1.0) * ledMaxLevel * 255.0)
		pixels[i].B = uint8(clamp(0, led.B, 1.0) * ledMaxLevel * 255.0)
	})

	leds.Write(sm, pixels)
}

const clockHz = 133000000

func blinkPinForever(sm *pio.PIOStateMachine, offset uint8, pin machine.Pin, freq uint) {
	blinkProgramInit(sm, offset, pin)
	sm.SetEnabled(true)

	fmt.Printf("Blinking pin %d at %d Hz\n", pin, freq)
	sm.Tx(uint32(clockHz / (2 * freq)))
}
