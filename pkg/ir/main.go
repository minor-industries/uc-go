package ir

import (
	"machine"
	"tinygo.org/x/drivers/irremote"
)

//const pinIRIn = machine.GP17
//const powerPin = machine.GP18

func Main(ch irremote.CommandHandler) {
	if powerPin != machine.NoPin {
		powerPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		powerPin.Set(true)
	}

	ir := irremote.NewReceiver(pinIRIn)
	ir.Configure()

	ir.SetCommandHandler(ch)
}
