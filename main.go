//go:build rp2040

package main

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"os"
	"time"
	"uc-go/app/bikelights"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/storage"
)

func main2() {
	rpcQueue := rpc.NewQueue(os.Stdout, 100)
	a := &bikelights.App{
		Logs: rpcQueue,
	}

	router := rpc.NewRouter()
	router.Register(map[string]rpc.Handler{
		"__sys__.dump-stored-logs": rpc.HandlerFunc(func(method string, body []byte) error {
			rpcQueue.Start()
			return nil
		}),
	})

	err := rfm69v2(a)
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "rfm69"))
	}

	router.Register(a.Handlers())
	go rpc.DecodeFrames(a.Logs, router)

	err = a.Run()
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "run exited with error"))
	} else {
		a.Logs.Log("run exited")
	}

	select {}
}

type Board struct{}

func (b Board) TxSPI(w, r []byte) error {
	return errors.New("not implemented")
}

func (b Board) Reset(b2 bool) error {
	return errors.New("not implemented")
}

func (b Board) WaitForD0Edge() {
	select {}
}

func rfm69v2(a *bikelights.App) error {
	board := &Board{}

	log := func(s string) {
		a.Logs.Log(s)
	}

	if err := rfm69.Run(board, log); err != nil {
		return errors.Wrap(err, "run rfm69")
	}

	return nil
}

func rfm69v1(a *bikelights.App) error {
	const REG_SYNCVALUE1 = 0x2F

	spi := machine.SPI0
	err := spi.Configure(machine.SPIConfig{
		Frequency: 64000,
		Mode:      machine.Mode3,
		SCK:       machine.GP2,
		SDO:       machine.GP3,
		SDI:       machine.GP4,
	})
	if err != nil {
		return errors.Wrap(err, "configure SPI")
	} else {
		a.Logs.Log("setup SPI")
	}

	rst := machine.GP6
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rst.Set(true)

	time.Sleep(300 * time.Millisecond) // TODO: shorten to optimal value

	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rst.Set(false)

	time.Sleep(300 * time.Millisecond) // TODO: shorten to optimal value

	CSn := machine.GP5
	CSn.Set(true)
	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CSn.Set(true)

	{
		for i := 0; i < 15; i++ {
			reg, err := readReg(spi, CSn, REG_SYNCVALUE1)
			if err != nil {
				return errors.Wrap(err, "read reg")
			}
			a.Logs.Log(fmt.Sprintf("val = 0x%02x, t=%s", reg, time.Now().String()))
			if reg == 0xAA {
				break
			}
			if err := writeReg(spi, CSn, REG_SYNCVALUE1, 0xAA); err != nil {
				return errors.Wrap(err, "write reg")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	return err
}

func readReg(spi *machine.SPI, CSn machine.Pin, addr byte) (byte, error) {
	CSn.Set(false)

	rx := make([]byte, 2)

	if err := spi.Tx(
		[]byte{addr & 0x7F, 0},
		rx,
	); err != nil {
		return 0, errors.Wrap(err, "tx")
	}

	CSn.Set(true)

	return rx[1], nil
}

func writeReg(spi *machine.SPI, CSn machine.Pin, addr byte, value byte) error {
	rx := make([]byte, 2)

	CSn.Set(false)

	if err := spi.Tx(
		[]byte{addr | 0x80, value},
		rx,
	); err != nil {
		return errors.Wrap(err, "tx")
	}

	CSn.Set(true)

	return nil
}

func resetConfig() {
	logs := rpc.NewQueue(os.Stdout, 100)

	lfs, err := storage.Setup(logs)
	if err != nil {
		panic("no")
	}

	if err := storage.WriteMsgp(logs, lfs, &cfg.DefaultConfig, "/cfg.msgp"); err != nil {
		panic("no")
	}

	select {}
}

func setConfig() {
	logs := rpc.NewQueue(os.Stdout, 100)

	lfs, err := storage.Setup(logs)
	if err != nil {
		panic("no")
	}

	c := cfg.Config{
		CurrentAnimation: "rainbow1",
		NumLeds:          150,
		StartIndex:       10,
		Length:           5,
		Scale:            0.5,
		MinScale:         0.04,
		ScaleIncr:        0.02,
	}

	if err := storage.WriteMsgp(logs, lfs, &c, "/cfg.msgp"); err != nil {
		panic("no")
	}
}

func main() {
	main2()
	//resetConfig()
	//setConfig()
}
