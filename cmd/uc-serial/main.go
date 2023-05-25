package main

import (
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/tarm/serial"
	"path/filepath"
	"time"
	"uc-go/cfg"
	"uc-go/protocol/framing"
	"uc-go/protocol/rpc"
	"uc-go/protocol/rpc/api"
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
			err := rpc.Send(device, "dump-stored-logs", nil)
			noErr(err)
		}
	}()

	go func() {
		err = framing.Decode(device, func(frame []byte) {
			//fmt.Printf("got frame: [%s]\n", frame)

			rpcMsg := &rpc.Request{}
			_, err := rpcMsg.UnmarshalMsg(frame)
			noErr(err)

			switch rpcMsg.Method {
			case "log":
				msg := &api.LogRequest{}
				_, err := msg.UnmarshalMsg(rpcMsg.Body)
				noErr(err)
				fmt.Println("got log:", msg.Message)

			case "show-config":
				fmt.Println("config:")
				fmt.Println(hex.Dump(rpcMsg.Body))
				msg := &cfg.Config{}
				_, err := msg.UnmarshalMsg(rpcMsg.Body)
				noErr(err)
				fmt.Println("config:", spew.Sdump(msg))

			default:
				fmt.Println("unknown message: " + rpcMsg.Method)
			}
		})
		noErr(err)
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			err = rpc.Send(device, "get-config", nil)
			noErr(err)
		}
	}()

	select {}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
