package ir

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/irremote"
)

const pinIRIn = machine.GP0

func Main() {
	ir := irremote.NewReceiver(pinIRIn)
	ir.Configure()

	ir.SetCommandHandler(func(data irremote.Data) {
		fmt.Printf("command: %d\r\n", data.Command)
	})

	for range time.NewTicker(5 * time.Second).C {
		fmt.Printf("running %s\r\n", time.Now().String())
	}
}
