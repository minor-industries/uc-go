package cfg

import (
	"sync"
	"uc-go/util"
)

type Config struct {
	CurrentAnimation string
	NumLeds          int
	StartIndex       int
	Length           float32

	Scale     float32
	MinScale  float32
	ScaleIncr float32
}

type SyncConfig struct {
	config Config
	lock   sync.Mutex
}

func NewSyncConfig(config Config) *SyncConfig {
	return &SyncConfig{config: config}
}

func (sc *SyncConfig) Edit(cb func(*Config)) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	cb(&sc.config)
}

func (sc *SyncConfig) SnapShot() Config {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	result := sc.config
	return result
}

func (sc *SyncConfig) ScaleUp() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	c := &sc.config
	c.Scale = util.Clamp(c.MinScale, c.Scale+c.ScaleIncr, 1.0)
}

func (sc *SyncConfig) ScaleDown() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	c := &sc.config
	c.Scale = util.Clamp(c.MinScale, c.Scale-c.ScaleIncr, 1.0)
}

func (sc *SyncConfig) SetAnimation(s string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	sc.config.CurrentAnimation = s
}
