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
	<-time.After(2 * time.Second)

	led := machine.A0
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.NEOPIXELS_POWER.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.NEOPIXELS_POWER.High()

	machine.NEOPIXELS.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neo := ws2812.New(machine.NEOPIXELS)
	if err := neo.WriteColors([]color.RGBA{{0, 0, 16, 0}}); err != nil {
		fmt.Println(errors.Wrap(err, "write neopixel"))
	}

	fmt.Println("should be active")

	bl := blikenlights.NewLight(func(on bool) {
		if on {
			neo.WriteColors([]color.RGBA{{0, 0, 16, 0}})
		} else {
			neo.WriteColors([]color.RGBA{{0, 16, 0, 0}})
		}
	})
	go bl.Run()

	bl.Seq([]int{2, 2})

	<-time.After(5 * time.Second)

	bl.Seq([]int{2, 10})

	select {}
}
