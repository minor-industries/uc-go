package main

import (
	"image/color"
	"machine"
	"math"
	"time"
	"tinygo.org/x/drivers/ws2812"
)

var (
	ledPin = machine.PA23
)

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.New(ledPin)

	t0 := time.Now()
	for t := range time.NewTicker(30 * time.Millisecond).C {
		dt := t.Sub(t0).Seconds()

		v := uint8(16*math.Sin(dt) + 16)

		leds.WriteColors([]color.RGBA{
			{v, 0, 0, 0},
			{0, v, 0, 0},
			{0, 0, v, 0},
		})
	}
}
