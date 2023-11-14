package main

import (
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"os"
	"path/filepath"
	"time"
)

func pollUsbSerial() string {
	polling := false
	for {
		dev, err := filepath.Glob("/dev/tty.usb*")
		noErr(err)

		switch len(dev) {
		case 1:
			fmt.Println("found", dev[0])
			return dev[0]
		case 0:
			if !polling {
				fmt.Println("polling for device")
			}
			polling = true
			time.Sleep(100 * time.Millisecond)
		default:
			panic(errors.New("found more than one serial device"))
		}
	}
}

func run() error {
	dev := pollUsbSerial()

	device, err := serial.OpenPort(&serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: 0,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	})
	noErr(err)

	_, err = io.Copy(os.Stdout, device)
	return err
}

func main() {
	for {
		err := run()
		fmt.Println("error:", err)
	}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
