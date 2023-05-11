//go:build rp2040

package main

import (
	"fmt"
	"image/color"
	"sync/atomic"
	"time"
	"tinygo/bounce"
	"tinygo/cfg"
	"tinygo/exe/ir"
	"tinygo/leds"
	"tinygo/pio"
	"tinygo/rainbow"
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

	tickDuration := 30 * time.Millisecond

	rb := rainbow.Rainbow1(
		&rainbow.App{Strip: strip},
		&rainbow.FaderConfig{TimeScale: 0.3},
	)

	count := uint32(0)
	t0 := time.Now()

	go func() {
		for range time.NewTicker(time.Second).C {
			count := atomic.LoadUint32(&count)
			dt := time.Now().Sub(t0)
			line := fmt.Sprintf(
				"count = %d, t=%s, fps=%0.02f",
				count,
				time.Now().String(),
				float64(count)/dt.Seconds(),
			)
			fmt.Printf(line + "\r\n")
		}
	}()

	for range time.NewTicker(tickDuration).C {
		atomic.AddUint32(&count, 1)

		if false {
			sim.Tick(0, tickDuration.Seconds())
		} else {
			t := float64(time.Now().UnixNano()) / 1e9
			rb(t, 0)
		}

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
