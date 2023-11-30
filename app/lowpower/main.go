package main

import (
	"device/arm"
	"fmt"
	"machine"
	"time"
	"uc-go/pkg/blikenlights"
)

var state = false

func callback(pin machine.Pin) {
	state = !state
	machine.LED.Set(state)
}

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bl := blikenlights.NewLight(machine.LED)
	go bl.Run()
	bl.Seq([]int{2, 2})

	<-time.After(10 * time.Second)
	bl.Seq([]int{2, 4})

	intr := machine.D1
	intr.Configure(machine.PinConfig{Mode: machine.PinInput})
	if err := intr.SetInterrupt(machine.PinRising, callback); err != nil {
		fmt.Println("error:", err.Error())
		bl.Seq([]int{2, 2, 2, 16})
	}

	for {
		sleep()
	}
}

func sleep() {
	//arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	arm.Asm("wfi")
}
