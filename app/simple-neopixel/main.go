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
	leds := ws2812.NewWS2812(ledPin)

	t0 := time.Now()
	buf := make([]color.RGBA, 150)

	for t := range time.NewTicker(30 * time.Millisecond).C {
		dt := t.Sub(t0).Seconds()

		v := uint8(16*math.Sin(dt) + 32)
		v = 0x05
		_ = v

		for i := range buf {
			buf[i].R = uint8(i)
			buf[i].G = 0
			buf[i].B = 0
		}

		leds.WriteColors(buf)
	}
}
