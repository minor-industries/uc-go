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

	i := 0
	t0 := time.Now()

	for {
		now := time.Now()
		dt := now.Sub(t0)
		sim.Tick(
			float64(now.UnixMicro()/1e6),
			dt.Seconds(),
		)

		writeColors(neo, colors, i)

		led.Low()
		time.Sleep(time.Millisecond * 500)

		led.High()
		time.Sleep(time.Millisecond * 500)

		i++
		i %= ledNum
	}
}

func writeColors(neo ws2812.Device, colors []color.RGBA, blank int) {
	for i := range colors {
		if i == blank {
			colors[i] = color.RGBA{
				R: 1,
				G: 0,
				B: 0,
				A: 255,
			}
		} else {
			colors[i] = color.RGBA{
				R: 32,
				G: 0,
				B: 0,
				A: 255,
			}
		}
	}

	err := neo.WriteColors(colors)
	if err != nil {
		fmt.Println("err: %s", err.Error())
	}
}
