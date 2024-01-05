package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.PA07
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		led.High()
		<-time.After(time.Second)
		led.Low()
		<-time.After(time.Second)
	}
}
