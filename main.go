//go:build rp2040

package main

import (
	"fmt"
	"image/color"
	"sync/atomic"
	"time"
	"tinygo.org/x/drivers/irremote"
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
	logs := make(chan irremote.Data, 10)

	ir.Main(func(data irremote.Data) {
		logs <- data
	})

	go func() {
		for log := range logs {
			line := fmt.Sprintf(
				"0x%02x, 0x%02x, 0x%02x 0x%02x",
				log.Code,
				log.Flags,
				log.Command,
				log.Address,
			)
			fmt.Println(line + "\r")
		}
	}()

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

	tickDuration := 30 * time.Millisecond

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

	animations := map[string]func(t, dt float64){
		"rainbow1": rainbow.Rainbow1(
			&rainbow.App{Strip: strip},
			&rainbow.FaderConfig{TimeScale: 0.3},
		),
		"rainbow2": rainbow.Rainbow2(
			&rainbow.App{Strip: strip},
			&rainbow.FaderConfig{TimeScale: 0.03},
		),
		"bounce": bounce.Bounce(
			&bounce.App{Strip: strip},
		).Tick,
	}

	currentAnimation := "bounce"

	f := func() {
		atomic.AddUint32(&count, 1)

		cb := animations[currentAnimation]
		t := float64(time.Now().UnixNano()) / 1e9
		cb(t, tickDuration.Seconds())
		writeColors(sm, pixels, strip)
	}

	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case <-ticker.C:
			f()
		}
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
