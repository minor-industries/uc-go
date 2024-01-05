package main

import (
	"fmt"
	"machine"
	"uc-go/pkg/blikenlights"
)

type logger struct{}

func (l *logger) Log(s string) {
	fmt.Println(s)
}

func (l *logger) Error(err error) {
	fmt.Printf("error: %v\n", err)
}

func (l *logger) Rpc(s string, i interface{}) error {
	fmt.Println("rpc: " + s)
	return nil
}

type Cfg struct {
	led machine.Pin
}

func setupLeds(cfg *Cfg) {
	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bl := blikenlights.NewLight(cfg.led)
	go bl.Run()
	bl.Seq([]int{2, 2})
}

func main() {
	cfg := &Cfg{
		led: machine.PA07,
	}

	setupLeds(cfg)

	select {}
}
