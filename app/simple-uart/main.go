package main

import (
	"fmt"
	"github.com/pkg/errors"
	"machine"
	"time"
)

func run() error {
	uart := machine.UART1
	if err := uart.Configure(machine.UARTConfig{
		BaudRate: 115200,
	}); err != nil {
		return errors.Wrap(err, "configure uart")
	}

	errCh := make(chan error)
	parts := make(chan string, 10)

	go func() {
		for {
			<-time.After(1 * time.Second)
			//fmt.Println("write")
			if _, err := uart.Write([]byte("\r\n")); err != nil {
				errCh <- errors.Wrap(err, "write")
				return
			}
		}
	}()

	go func() {
		var buf [64]byte
		for {
			n, err := uart.Read(buf[:])
			if err != nil {
				errCh <- errors.Wrap(err, "read")
				return
			}
			if n > 0 {
				parts <- string(buf[:n])
				continue
			}
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		for part := range parts {
			fmt.Printf(part)
		}
	}()

	select {}
}

func main() {
	err := run()
	for {
		fmt.Println("error:", err)
		<-time.After(time.Second)
	}
}
