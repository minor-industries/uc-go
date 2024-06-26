// Code generated by pioasm; DO NOT EDIT.

//go:build rp2040
// +build rp2040

package leds

import (
	"github.com/minor-industries/uc-go/pkg/pio"
	"machine"
)

// ws2812

const ws2812WrapTarget = 0
const ws2812Wrap = 3

const ws2812T1 = 2
const ws2812T2 = 5
const ws2812T3 = 3

var ws2812Program = pio.PIOProgram{
	Instructions: []uint16{
		//     .wrap_target
		0x6221, //  0: out    x, 1            side 0 [2]
		0x1123, //  1: jmp    !x, 3           side 1 [1]
		0x1400, //  2: jmp    0               side 1 [4]
		0xa442, //  3: nop                    side 0 [4]
		//     .wrap
	},
	Origin: -1,
}

func ws2812ProgramDefaultConfig(offset uint8) pio.PIOStateMachineConfig {
	cfg := pio.DefaultStateMachineConfig()
	cfg.SetWrap(offset+ws2812WrapTarget, offset+ws2812Wrap)
	cfg.SetSideSet(1, false, false)
	return cfg
}

const freq = 125_000_000

func ws2812ProgramInit(
	sm *pio.PIOStateMachine,
	offset uint8,
	pin machine.Pin,
) {
	pin.Configure(machine.PinConfig{Mode: machine.PinPIO0})
	sm.SetConsecutivePinDirs(pin, 1, true)
	cfg := ws2812ProgramDefaultConfig(offset)
	cfg.SetSetPins(pin, 1)
	cfg.SetOutShift(false, true, 24)

	// Frequency = clock freq / (CLKDIV_INT + CLKDIV_FRAC / 256)

	// TODO: set fifo join
	//cycles_per_bit := ws2812T1 + ws2812T2 + ws2812T3

	cfg.SetClkDivIntFrac(15, 159)
	sm.Init(offset, &cfg)
}
