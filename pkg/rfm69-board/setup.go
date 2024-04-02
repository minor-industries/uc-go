package rfm69_board

import (
	"github.com/minor-industries/rfm69"
	"github.com/minor-industries/uc-go/pkg/logger"
	"github.com/minor-industries/uc-go/pkg/rfm69-board/cfg"
	"github.com/minor-industries/uc-go/pkg/spi"
	"github.com/minor-industries/uc-go/pkg/storage"
	"github.com/minor-industries/uc-go/pkg/util"
	"github.com/pkg/errors"
	"machine"
	"tinygo.org/x/tinyfs/littlefs"
)

type PinCfg struct {
	Rst  machine.Pin
	Intr machine.Pin
}

func SetupRfm69(
	env *cfg.Config,
	Spi *spi.SPI,
	cfg *PinCfg,
	log func(s string),
) (*rfm69.Radio, error) {
	rst := cfg.Rst
	rst.Configure(machine.PinConfig{Mode: machine.PinOutput})

	board, err := NewBoard(
		Spi,
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
	initialNodeAddr = 0xd0
)

func LoadConfig(
	logs logger.Logger,
	lfs *littlefs.LFS,
	forceReInit bool,
) (
	*util.SyncConfig[cfg.Config],
	error,
) {
	config, err := storage.LoadConfig[*cfg.Config](
		lfs,
		logs,
		configFile,
		forceReInit,
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
