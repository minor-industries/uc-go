package bbq

import (
	"github.com/pkg/errors"
	"machine"
	"math/rand"
	"time"
	"uc-go/pkg/protocol/rpc"
	rfm69_board "uc-go/pkg/rfm69-board"
	cfg3 "uc-go/pkg/rfm69-board/cfg"
)

const dstAddr = 2

func Run(logs *rpc.Queue) error {
	env, err := rfm69_board.LoadConfig(logs)
	if err != nil {
		return errors.Wrap(err, "load config")
	}

	cfg.led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	stopLeds := make(chan struct{})
	go ledControl(stopLeds)
	go func() {
		<-time.After(5 * time.Second)
		close(stopLeds)
	}()

	log := func(s string) {
		logs.Log(s)
	}

	envSnapshot := env.SnapShot()

	err = mainLoop(
		logs,
		&envSnapshot,
		log,
	)
	if err != nil {
		return errors.Wrap(err, "mainloop")
	}

	return errors.New("run exited")
}

func ledControl(done <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	val := true

	for {
		select {
		case <-ticker.C:
			cfg.led.Set(val)
			val = !val
		case <-done:
			return
		}
	}
}

func mainLoop(
	logs *rpc.Queue,
	env *cfg3.Config,
	log func(s string),
) error {
	randSource := rand.New(rand.NewSource(int64(env.NodeAddr)))

	_, err := rfm69_board.SetupRfm69(env, &cfg.Rfm, log)
	if err != nil {
		return errors.Wrap(err, "rfm69")
	}

	readAndSend := func() error {
		return errors.New("read and send not implemented")
	}

	for {
		if err := readAndSend(); err != nil {
			logs.Error(err)
		}

		sleep := time.Duration(4000+randSource.Intn(2000)) * time.Millisecond
		time.Sleep(sleep)
	}
}
