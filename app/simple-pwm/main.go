package main

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"math"
	"time"
)

var period uint64 = 1e9 / 100_000

func run() error {
	pwm := machine.TCC0

	channel, err := pwm.Channel(machine.D2)
	if err != nil {
		return errors.Wrap(err, "channel")
	}

	err = pwm.Configure(machine.PWMConfig{Period: period})
	if err != nil {
		return errors.Wrap(err, "configure pwm")
	}

	t0 := time.Now()
	for {
		t := time.Now().Sub(t0)
		v := 0.5 + 0.5*math.Sin(0.25*t.Seconds())
		vf := float64(pwm.Top()) * v
		vi := uint32(vf)
		fmt.Println(t, v, vi)
		pwm.Set(channel, vi)
	}

	fmt.Println("channel", channel)
	return nil
}

func main() {
	for {
		err := run()
		if err != nil {
			fmt.Println("error:", err)
		}
		<-time.After(5 * time.Second)
	}
}
