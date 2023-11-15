package rfm69_board

import (
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"machine"
	"uc-go/pkg/logger"
	"uc-go/pkg/rfm69-board/cfg"
	"uc-go/pkg/spi"
	"uc-go/pkg/storage"
	"uc-go/pkg/util"
)

type PinCfg struct {
	// rfm
	Spi *spi.Config

	Rst  machine.Pin
	Intr machine.Pin
}

func SetupRfm69(
	env *cfg.Config,
	cfg *PinCfg,
	spi *spi.SPI,
	log func(s string),
) (*rfm69.Radio, error) {
	rst := cfg.Rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	board, err := NewBoard(
		spi,
		rst,
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

func LoadConfig(logs logger.Logger) (
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
