package main

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/irremote"
)

func main() {
	led := machine.PA23
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	value := true

	led.Set(value)

	irPin := machine.PA15

	ir := irremote.NewReceiver(irPin)
	ir.Configure()

	ch := make(chan irremote.Data, 10)

	ir.SetCommandHandler(func(data irremote.Data) {
		ch <- data
	})

	lastUpdate := time.Now()
	for data := range ch {
		now := time.Now()
		elapsed := now.Sub(lastUpdate)
		lastUpdate = now

		fmt.Println("IR", data.Address, data.Code, elapsed)

		if elapsed > time.Second {
			value = !value
			led.Set(value)
		}
	}
}
