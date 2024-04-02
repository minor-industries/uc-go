package main

import (
	"device/arm"
	"device/sam"
	"github.com/minor-industries/uc-go/pkg/blikenlights"
	"machine"
	"time"
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
	bl.Off()

	intr := machine.D1
	intr.Configure(machine.PinConfig{Mode: machine.PinInput})
	intr.SetInterrupt(machine.PinRising, callback)

	configGCLK6()
	//
	//sam.PM.SetSLEEP_IDLE()
	//sam.EIC.WAKEUP.Set(0x00000000)
	//sam.EIC.WAKEUP.ClearBits(1 << intr)
	sam.EIC.WAKEUP.SetBits(1 << intr)
	for sam.EIC.STATUS.HasBits(sam.EIC_STATUS_SYNCBUSY) {
	}

	//interrupt.New(sam.IRQ_RTC, func(i interrupt.Interrupt) {
	//})

	sam.PM.APBAMASK.SetBits(sam.PM_APBAMASK_GCLK_ |
		sam.PM_APBAMASK_EIC_ |
		sam.PM_APBAMASK_RTC_)

	machine.LED.Low()
	count := 0
	for {
		sleep()
		count++
		if count%100 == 0 {
			state = !state
			machine.LED.Set(state)
		}
		//time.Sleep(50 * time.Millisecond)
	}
}

func configGCLK6() {
	sam.GCLK.CLKCTRL.ClearBits(sam.GCLK_CLKCTRL_CLKEN)
	waitForSync()

	//sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_EIC<<sam.GCLK_CLKCTRL_ID_Pos |
	//	sam.GCLK_CLKCTRL_GEN_GCLK0<<sam.GCLK_CLKCTRL_GEN_Pos |
	//	sam.GCLK_CLKCTRL_CLKEN)
	//waitForSync()

	// *****************
	sam.GCLK.GENDIV.Set(6 << sam.GCLK_GENDIV_ID_Pos)
	waitForSync()

	sam.GCLK.GENCTRL.Set((6 << sam.GCLK_GENCTRL_ID_Pos) |
		(sam.GCLK_GENCTRL_SRC_OSCULP32K << sam.GCLK_GENCTRL_SRC_Pos) |
		sam.GCLK_GENCTRL_RUNSTDBY | // this one is always on
		sam.GCLK_GENCTRL_GENEN)
	waitForSync()

	sam.GCLK.CLKCTRL.Set((sam.GCLK_CLKCTRL_ID_EIC << sam.GCLK_CLKCTRL_ID_Pos) |
		(sam.GCLK_CLKCTRL_GEN_GCLK6 << sam.GCLK_CLKCTRL_GEN_Pos) |
		sam.GCLK_CLKCTRL_CLKEN)
	waitForSync()
	// *****************

	sam.GCLK.GENDIV.Set(2 << sam.GCLK_GENDIV_ID_Pos)
	waitForSync()

	sam.GCLK.GENCTRL.Set((2 << sam.GCLK_GENCTRL_ID_Pos) |
		(sam.GCLK_GENCTRL_SRC_OSC32K << sam.GCLK_GENCTRL_SRC_Pos) |
		sam.GCLK_GENCTRL_GENEN)
	waitForSync()

	sam.GCLK.GENCTRL.SetBits(sam.GCLK_GENCTRL_RUNSTDBY)
	waitForSync()

	//GCLK->GENCTRL.bit.RUNSTDBY = 1;  //GCLK6 run standby
	//sam.GCLK_GENCTRL_RUNSTDBY |

	// Use GCLK2 for RTC
	sam.GCLK.CLKCTRL.Set((sam.GCLK_CLKCTRL_ID_RTC << sam.GCLK_CLKCTRL_ID_Pos) |
		(sam.GCLK_CLKCTRL_GEN_GCLK2 << sam.GCLK_CLKCTRL_GEN_Pos) |
		sam.GCLK_CLKCTRL_CLKEN)
	waitForSync()

	//sam.GCLK.CLKCTRL.Set(sam.GCLK_CLKCTRL_ID_EIC<<sam.GCLK_CLKCTRL_ID_Pos |
	//	sam.GCLK_CLKCTRL_GEN_GCLK0<<sam.GCLK_CLKCTRL_GEN_Pos |
	//	sam.GCLK_CLKCTRL_CLKEN)
	//waitForSync()

	//sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_CLKEN | sam.GCLK_CLKCTRL_GEN_GCLK6 | sam.GCLK_CLKCTRL_ID_EIC)
	//waitForSync()
	//}

	//sam.GCLK.GENCTRL.Set(sam.GCLK_GENCTRL_RUNSTDBY | sam.GCLK_GENCTRL_)
	//waitForSync()
	//}

	sam.GCLK.CLKCTRL.SetBits(sam.GCLK_CLKCTRL_CLKEN)
	waitForSync()
	//
	sam.NVMCTRL.CTRLB.SetBits(sam.NVMCTRL_CTRLB_SLEEPPRM_DISABLED)
}

func waitForSync() {
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
	arm.SCB.SCR.SetBits(arm.SCB_SCR_SLEEPDEEP)
	arm.Asm("wfi")
}
