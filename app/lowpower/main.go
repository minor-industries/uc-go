package main

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
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

	<-time.After(5 * time.Second)
	bl.Seq([]int{2, 4})

	intr := machine.D1
	intr.Configure(machine.PinConfig{Mode: machine.PinInput})
	if err := intr.SetInterrupt(machine.PinRising, callback); err != nil {
		fmt.Println("error:", err.Error())
		bl.Seq([]int{2, 2, 2, 16})
	}

	configGCLK6()
	//
	//sam.PM.SetSLEEP_IDLE()
	//sam.EIC.WAKEUP.Set(0x00000000)
	//sam.EIC.WAKEUP.ClearBits(1 << intr)
	sam.EIC.WAKEUP.SetBits(1 << intr)
	for sam.EIC.STATUS.HasBits(sam.EIC_STATUS_SYNCBUSY) {
	}

	interrupt.New(sam.IRQ_RTC, func(i interrupt.Interrupt) {
		state = !state
		machine.LED.Set(state)
	})

	sam.PM.APBAMASK.SetBits(sam.PM_APBAMASK_GCLK_ | sam.PM_APBAMASK_EIC_)

	for {
		sleep()
	}
}

func configGCLK6() {
	sam.GCLK.CLKCTRL.ClearBits(sam.GCLK_CLKCTRL_CLKEN)
	waitClockSync()

	//sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_EIC<<sam.GCLK_CLKCTRL_ID_Pos |
	//	sam.GCLK_CLKCTRL_GEN_GCLK0<<sam.GCLK_CLKCTRL_GEN_Pos |
	//	sam.GCLK_CLKCTRL_CLKEN)
	//waitClockSync()

	// *****************
	sam.GCLK.GENDIV.Set(6 << sam.GCLK_GENDIV_ID_Pos)
	waitClockSync()

	sam.GCLK.GENCTRL.Set((6 << sam.GCLK_GENCTRL_ID_Pos) |
		(sam.GCLK_GENCTRL_SRC_OSCULP32K << sam.GCLK_GENCTRL_SRC_Pos) |
		sam.GCLK_GENCTRL_RUNSTDBY |
		sam.GCLK_GENCTRL_GENEN)
	waitClockSync()

	sam.GCLK.CLKCTRL.Set((sam.GCLK_CLKCTRL_ID_EIC << sam.GCLK_CLKCTRL_ID_Pos) |
		(sam.GCLK_CLKCTRL_GEN_GCLK6 << sam.GCLK_CLKCTRL_GEN_Pos) |
		sam.GCLK_CLKCTRL_CLKEN)
	waitClockSync()
	// *****************

	//sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_EIC<<sam.GCLK_CLKCTRL_ID_Pos |
	//	sam.GCLK_CLKCTRL_GEN_GCLK0<<sam.GCLK_CLKCTRL_GEN_Pos |
	//	sam.GCLK_CLKCTRL_CLKEN)
	//waitClockSync()

	//sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_CLKEN | sam.GCLK_CLKCTRL_GEN_GCLK6 | sam.GCLK_CLKCTRL_ID_EIC)
	//waitClockSync()
	//}

	//sam.GCLK.GENCTRL.Set(sam.GCLK_GENCTRL_RUNSTDBY | sam.GCLK_GENCTRL_)
	//waitClockSync()
	//}

	sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_CLKEN)
	waitClockSync()
	//
	//sam.NVMCTRL.CTRLB.SetBits(sam.NVMCTRL_CTRLB_SLEEPPRM_DISABLED)
}

func waitClockSync() {
	for sam.GCLK.STATUS.HasBits(sam.GCLK_STATUS_SYNCBUSY) {
	}
}

/*
static void configGCLK6()
{
	// enable EIC clock
	GCLK->CLKCTRL.bit.CLKEN = 0; //disable GCLK module
	while (GCLK->STATUS.bit.SYNCBUSY);

	GCLK->CLKCTRL.reg = (uint16_t) (GCLK_CLKCTRL_CLKEN | GCLK_CLKCTRL_GEN_GCLK6 | GCLK_CLKCTRL_ID( GCM_EIC )) ;  //EIC clock switched on GCLK6
	while (GCLK->STATUS.bit.SYNCBUSY);

	GCLK->GENCTRL.reg = (GCLK_GENCTRL_GENEN | GCLK_GENCTRL_SRC_OSCULP32K | GCLK_GENCTRL_ID(6));  //source for GCLK6 is OSCULP32K
	while (GCLK->STATUS.reg & GCLK_STATUS_SYNCBUSY);

	GCLK->GENCTRL.bit.RUNSTDBY = 1;  //GCLK6 run standby
	while (GCLK->STATUS.reg & GCLK_STATUS_SYNCBUSY);

	Errata: Make sure that the Flash does not power all the way down
     	* when in sleep mode.

	NVMCTRL->CTRLB.bit.SLEEPPRM = NVMCTRL_CTRLB_SLEEPPRM_DISABLED_Val;
}
*/

func sleep() {
	time.Sleep(100 * time.Millisecond)
	//arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	//arm.Asm("wfi")
}
