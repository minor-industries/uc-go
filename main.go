package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/ws2812"
	"tinygo/bounce"
	"tinygo/cfg"
	"tinygo/strip"
)

const (
	ledPin      = machine.GP0 // NeoPixels pin
	ledNum      = 150         // number of NeoPixels
	ledMaxLevel = 0.5         // brightness level of NeoPxels (0~1)
)

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neo := ws2812.New(ledPin)

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	strip := strip.NewStrip(&cfg.Cfg{
		NumLeds:    150,
		StartIndex: 0,
		Length:     5.0,
	})
	sim := bounce.Bounce(&bounce.App{Strip: strip})

	tick := 30 * time.Millisecond
	frames := 0

	t0 := time.Now()

	for range time.NewTicker(tick).C {
		//now := time.Now()
		//ts := now.Unix() / int64(time.Second)
		led.Set(!led.Get())

		sim.Tick(
			0,
			tick.Seconds(),
		)

		writeColors(neo, strip)
		frames++

		if frames%100 == 0 {
			now := time.Now()
			dt := now.Sub(t0)
			fmt.Println(frames, dt, float64(frames)/dt.Seconds(), "\r")
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

func writeColors(neo ws2812.Device, st *strip.Strip) {
	var colors [ledNum]color.RGBA

	st.Each(func(i int, led *strip.Led) {
		colors[i].R = uint8(clamp(0, led.R, 1.0) * ledMaxLevel * 255.0)
		colors[i].G = uint8(clamp(0, led.G, 1.0) * ledMaxLevel * 255.0)
		colors[i].B = uint8(clamp(0, led.B, 1.0) * ledMaxLevel * 255.0)
	})

	err := neo.WriteColors(colors[:])
	if err != nil {
		fmt.Println("err: %s", err.Error())
	}
}
