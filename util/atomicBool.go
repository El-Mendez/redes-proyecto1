package utils

import (
	"sync"
)

type AtomicBool struct {
	mutex sync.RWMutex
	value bool
}

func (atomicBool *AtomicBool) Get() bool {
	atomicBool.mutex.RLock()
	defer atomicBool.mutex.RUnlock()

	return atomicBool.value
}

func (atomicBool *AtomicBool) Set(v bool) {
	atomicBool.mutex.Lock()
	defer atomicBool.mutex.Unlock()

	atomicBool.value = v
}
