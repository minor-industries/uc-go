package spi

import (
	"github.com/pkg/errors"
	"machine"
	"sync"
)

type Config struct {
	Spi    *machine.SPI
	Config *machine.SPIConfig
	Cs     machine.Pin
}

type SPI struct {
	lock *sync.Mutex
	cfg  *Config
}

func NewSPI(
	config *Config,
	lock *sync.Mutex,
) *SPI {
	if config.Cs != machine.NoPin {
		config.Cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
		config.Cs.Set(true)
	}

	s := &SPI{
		cfg:  config,
		lock: lock,
	}

	return s
}

func (s *SPI) TxSPI(w, r []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.cfg.Spi.Configure(*s.cfg.Config); err != nil {
		return errors.Wrap(err, "configure spi")
	}

	if s.cfg.Cs != machine.NoPin {
		s.cfg.Cs.Low()
	}

	err := s.cfg.Spi.Tx(w, r)

	if s.cfg.Cs != machine.NoPin {
		s.cfg.Cs.High()
	}

	return errors.Wrap(err, "tx spi")
}
