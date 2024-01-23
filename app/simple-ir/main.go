package main

import (
	"machine"
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

	ir.SetCommandHandler(func(data irremote.Data) {
		value = !value
		led.Set(value)
	})

	select {}
}
