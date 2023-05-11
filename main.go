//go:build rp2040

package main

import (
	"image/color"
	"time"
	"tinygo/bounce"
	"tinygo/cfg"
	"tinygo/exe/ir"
	"tinygo/leds"
	"tinygo/pio"
	"tinygo/strip"
)

const (
	ledMaxLevel = 0.5 // brightness level of NeoPxels (0~1)
)

func main() {
	ir.Main()
	////pioMain()
	sm := leds.Setup()

	runLeds(sm)
}

func runLeds(sm *pio.PIOStateMachine) {
	pixels := make([]color.RGBA, 150)

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
