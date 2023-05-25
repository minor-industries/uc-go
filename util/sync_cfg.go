package util

import (
	"sync"
)

type SyncConfig[T any] struct {
	config T
	lock   sync.Mutex
}

func NewSyncConfig[T any](config T) *SyncConfig[T] {
	return &SyncConfig[T]{config: config}
}

func (sc *SyncConfig[T]) Edit(cb func(*T)) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	cb(&sc.config)
}

func (sc *SyncConfig[T]) SnapShot() T {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	result := sc.config
	return result
}
