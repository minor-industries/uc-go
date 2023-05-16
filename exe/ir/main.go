package ir

import (
	"machine"
	"tinygo.org/x/drivers/irremote"
)

//const pinIRIn = machine.GP17
//const powerPin = machine.GP18

const pinIRIn = machine.GP21 // TODO: should be config?
const powerPin = machine.GP22

func Main(ch irremote.CommandHandler) {
	powerPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	powerPin.Set(true)

	ir := irremote.NewReceiver(pinIRIn)
	ir.Configure()

	ir.SetCommandHandler(ch)
}
