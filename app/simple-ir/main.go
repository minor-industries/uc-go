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

	ch := make(chan irremote.Data, 10)

	ir.SetCommandHandler(func(data irremote.Data) {
		ch <- data
	})

	for data := range ch {
		//fmt.Println("IR", data.Address, data.Code, data.Command)
		switch data.Command {
		case 16:
			led.Low()
		case 17:
			led.High()
		}
	}
}
