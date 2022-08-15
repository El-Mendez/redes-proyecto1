package utils

import (
	"sync"
)

// AtomicBool functions like a thread-safe boolean.
type AtomicBool struct {
	mutex sync.RWMutex
	value bool
}

// Get returns the value of the boolean with thread-safety.
func (atomicBool *AtomicBool) Get() bool {
	atomicBool.mutex.RLock()
	defer atomicBool.mutex.RUnlock()

	return atomicBool.value
}

// Set changes the value of the boolean with thread-safety.
func (atomicBool *AtomicBool) Set(v bool) {
	atomicBool.mutex.Lock()
	defer atomicBool.mutex.Unlock()

	atomicBool.value = v
}
