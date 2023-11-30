package main

import (
	"fmt"
	"github.com/pkg/errors"
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/ws2812"
	"uc-go/pkg/blikenlights"
)

func main() {
	led := machine.A0
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.NEOPIXELS_POWER.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXELS_POWER.High()

	machine.NEOPIXELS.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neo := ws2812.New(machine.NEOPIXELS)
	if err := neo.WriteColors([]color.RGBA{{0, 0, 16, 0}}); err != nil {
		fmt.Println(errors.Wrap(err, "write neopixel"))
	}

	bl := blikenlights.NewLight(blikenlights.BlinkerFunc(func(on bool) {
		if on {
			neo.WriteColors([]color.RGBA{{0, 0, 16, 0}})
		} else {
			neo.WriteColors([]color.RGBA{{0, 0, 0, 0}})
		}
	}))
	go bl.Run()

	bl.Off()

	<-time.After(5 * time.Second)

	bl.On()

	<-time.After(5 * time.Second)

	bl.Seq([]int{4, 4, 4, 16})

	select {}
}
