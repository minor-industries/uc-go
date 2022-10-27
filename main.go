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
	ledNum      = 13          // number of NeoPixels
	ledMaxLevel = 0.5         // brightness level of NeoPxels (0~1)
)

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neo := ws2812.New(ledPin)
	colors := make([]color.RGBA, ledNum)

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	sim := bounce.Bounce(&bounce.App{Strip: strip.NewStrip(&cfg.Cfg{
		NumLeds:    150,
		StartIndex: 0,
		Length:     5.0,
	})})

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

		writeColors(neo, colors)
		frames++

		if frames%100 == 0 {
			now := time.Now()
			dt := now.Sub(t0)
			fmt.Println(frames, dt, float64(frames)/dt.Seconds(), "\r")
		}
	}
}

func writeColors(neo ws2812.Device, colors []color.RGBA) {
	for i := range colors {
		colors[i] = color.RGBA{
			R: 32,
			G: 0,
			B: 0,
			A: 255,
		}
	}

	err := neo.WriteColors(colors)
	if err != nil {
		fmt.Println("err: %s", err.Error())
	}
}
