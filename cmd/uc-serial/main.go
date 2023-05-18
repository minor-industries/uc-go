package main

import (
	"fmt"
	"github.com/tarm/serial"
	"path/filepath"
	"time"
	"uc-go/protocol/framing"
)

func main() {
	dev, err := filepath.Glob("/dev/tty.usb*")
	noErr(err)

	if len(dev) != 1 {
		panic(fmt.Errorf("found %d serial devices", len(dev)))
	}

	device, err := serial.OpenPort(&serial.Config{
		Name:        dev[0],
		Baud:        115200,
		ReadTimeout: 0,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	})

	go func() {
		for range time.NewTicker(time.Second).C {
			msg := "hello from " + time.Now().String()
			_, err := device.Write(framing.Encode([]byte(msg)))
			noErr(err)
		}
	}()

	err = framing.Decode(device, func(frame []byte) {
		fmt.Printf("got frame: [%s]\n", frame)
	})
	noErr(err)
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
