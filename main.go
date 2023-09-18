//go:build rp2040

package main

import (
	"fmt"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"os"
	"time"
	"tinygo.org/x/drivers/aht20"
	"uc-go/app/bikelights"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
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

	i2c := machine.I2C0

	err := i2c.Configure(machine.I2CConfig{
		SDA: machine.GP0,
		SCL: machine.GP1,
	})
	if err != nil {
		a.Logs.Error(errors.Wrap(err, "configure i2c"))
	} else {
		go readTemp(i2c, a.Logs)
	}

	go func() {
		err := rfm69v2(a)
		if err != nil {
			a.Logs.Error(errors.Wrap(err, "rfm69"))
		}
	}()

	//router.Register(a.Handlers())
	go rpc.DecodeFrames(a.Logs, router)

	//err = a.Run()
	//if err != nil {
	//	a.Logs.Error(errors.Wrap(err, "run exited with error"))
	//} else {
	//	a.Logs.Log("run exited")
	//}

	select {}
}

func readTemp(i2c *machine.I2C, logs *rpc.Queue) {
	sensor := aht20.New(i2c)
	sensor.Configure()
	sensor.Reset()

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err := sensor.Read()
		if err != nil {
			logs.Error(errors.Wrap(err, "read sensor"))
			continue
		}
		t := sensor.Celsius()
		logs.Log(fmt.Sprintf("temperature = %0.01fC %0.01fF", t, (t*9/5)+32))
	}
}

func rfm69v2(a *bikelights.App) error {
	log := func(s string) {
		a.Logs.Log(s)
	}

	rst := machine.GP6
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := machine.SPI0
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode3,
		SCK:  machine.GP2,
		SDO:  machine.GP3,
		SDI:  machine.GP4,
	})
	if err != nil {
		return errors.Wrap(err, "configure SPI")
	} else {
		a.Logs.Log("setup SPI")
	}

	CSn := machine.GP5
	CSn.Set(true)
	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CSn.Set(true)

	board, err := rfm69_board.NewBoard(
		spi,
		rst,
		CSn,
		machine.GP7,
		log,
	)
	if err != nil {
		return errors.Wrap(err, "new board")
	}

	if err := rfm69.Setup(board, log); err != nil {
		return errors.Wrap(err, "run rfm69")
	}

	if err := rfm69.Rx(board, log); err != nil {
		return errors.Wrap(err, "rx rfm69")
	}

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
