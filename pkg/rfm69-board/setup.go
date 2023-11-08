package rfm69_board

import (
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"sync"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/rfm69-board/cfg"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
)

type PinCfg struct {
	// rfm
	Spi *machine.SPI

	Rst  machine.Pin
	Intr machine.Pin

	Sck machine.Pin
	Sdo machine.Pin
	Sdi machine.Pin
	Csn machine.Pin
}

func SetupRfm69(
	env *cfg.Config,
	cfg *PinCfg,
	spiLock *sync.Mutex,
	log func(s string),
) (*rfm69.Radio, error) {
	rst := cfg.Rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	spi := cfg.Spi
	err := spi.Configure(machine.SPIConfig{
		Mode: machine.Mode0,
		SCK:  cfg.Sck,
		SDO:  cfg.Sdo,
		SDI:  cfg.Sdi,
	})

	log("hello")

	//return nil, err

	if err != nil {
		return nil, errors.Wrap(err, "configure SPI")
	} else {
		log("setup SPI")
	}

	CSn := cfg.Csn

	CSn.Set(true)
	CSn.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CSn.Set(true)

	board, err := NewBoard(
		spi,
		spiLock,
		rst,
		CSn,
		cfg.Intr,
		log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new board")
	}

	radio := rfm69.NewRadio(board, log, env.NodeAddr, env.TxPower)

	if err := radio.Setup(rfm69.RF69_433MHZ); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	return radio, nil
}

const (
	configFile = "/radio-cfg.msgp"

	initialTxPower  = 20
	initialNodeAddr = 0xee
)

func LoadConfig(logs *rpc.Queue) (
	*util.SyncConfig[cfg.Config],
	error,
) {
	lfs, err := storage.Setup(logs)
	if err != nil {
		return nil, errors.Wrap(err, "setup storage")
	}

	if lfs == nil {
		return nil, errors.New("no lfs")
	}

	config, err := storage.LoadConfig[*cfg.Config](
		lfs,
		logs,
		configFile,
		&cfg.Config{
			NodeAddr: initialNodeAddr,
			TxPower:  initialTxPower,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	env := util.NewSyncConfig(*config)
	return env, nil
}
