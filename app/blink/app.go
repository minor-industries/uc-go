package blink

import (
	"machine"
	"time"
	"uc-go/pkg/protocol/rpc"
)

const dstAddr = 2

func Run(logs *rpc.Queue) error {
	go func() {
		for {
			logs.Log(time.Now().String())
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for {
		machine.LED.High()
		time.Sleep(50 * time.Millisecond)
		machine.LED.Low()
		time.Sleep(50 * time.Millisecond)
	}
}
