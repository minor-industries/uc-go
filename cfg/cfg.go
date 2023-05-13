package cfg

import (
	"sync"
	"tinygo/util"
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
	Config Config
	lock   sync.Mutex
}

func (sc *SyncConfig) Edit(cb func(*Config)) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	cb(&sc.Config)
}

func (sc *SyncConfig) SnapShot() Config {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	result := sc.Config
	return result
}

func (sc *SyncConfig) ScaleUp() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	c := sc.Config
	c.Scale = util.Clamp(c.MinScale, c.Scale+c.ScaleIncr, 1.0)
}

func (sc *SyncConfig) ScaleDown() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	c := sc.Config
	c.Scale = util.Clamp(c.MinScale, c.Scale-c.ScaleIncr, 1.0)
}

func (sc *SyncConfig) SetAnimation(s string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	sc.Config.CurrentAnimation = s
}
