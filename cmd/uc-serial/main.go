package main

import (
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"path/filepath"
	"time"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/protocol/framing"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/protocol/rpc/api"
)

var opts struct {
	Remove     bool `long:"remove" optional:"true"`
	SetConfig  bool `long:"set-config" optional:"true"`
	ShowConfig bool `long:"show-config" optional:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	noErr(err)

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
			err := rpc.Send(device, "__sys__.dump-stored-logs", nil)
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

	if opts.Remove {
		go func() {
			time.Sleep(2 * time.Second)
			err = rpc.Send(device, "reset-config", nil)
			noErr(err)
		}()
	}

	if opts.SetConfig {
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("setting config")
			msg := &cfg.Config{
				CurrentAnimation: "bounce",
				NumLeds:          10,
				StartIndex:       0,
				Length:           5.0,
				Scale:            0.5,
				MinScale:         0.04,
				ScaleIncr:        0.02,
			}
			err = rpc.Send(device, "set-config", msg)
			noErr(err)
		}()
	}

	if opts.ShowConfig {
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("showing config")
			err = rpc.Send(device, "get-config", nil)
			noErr(err)
		}()
	}

	select {}
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}
