package cfg

import "sync"

type Config struct {
	CurrentAnimation string
	NumLeds          int
	StartIndex       int
	Length           float64
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
