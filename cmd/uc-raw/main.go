package main

import (
	"fmt"
	"github.com/tarm/serial"
	"io"
	"os"
	"path/filepath"
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
		io.Copy(os.Stdout, device)
	}()

	select {}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
