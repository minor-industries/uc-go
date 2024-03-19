package main

import (
	"fmt"
	"machine"
	"tinygo.org/x/drivers/irremote"
)

const irPin = machine.PA15

//const irPin = machine.D5

func main() {
	ir := irremote.NewReceiver(irPin)
	ir.Configure()

	ch := make(chan irremote.Data, 10)

	ir.SetCommandHandler(func(data irremote.Data) {
		ch <- data
	})

	irCount := 0

	for {
		select {
		case data := <-ch:
			if data.Flags&irremote.DataFlagIsRepeat != 0 {
				continue
			}
			switch data.Command {
			case 16:
				irCount++
			case 17:
				irCount--
			}

			fmt.Println(irCount)
		}
	}
}
