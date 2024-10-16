package slocks

import "sync"

// CHANGE(limechain): Provides a per-slot mutex locker.

type PerSlotLocker struct {
	mu    sync.Mutex
	locks map[uint64]*sync.Mutex
}

// lock returns the lock of the given slot.
func (l *PerSlotLocker) lock(slot uint64) *sync.Mutex {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.locks == nil {
		l.locks = make(map[uint64]*sync.Mutex)
	}
	if _, ok := l.locks[slot]; !ok {
		l.locks[slot] = new(sync.Mutex)
	}
	return l.locks[slot]
}

// Lock locks the slot's mutex.
func (l *PerSlotLocker) Lock(slot uint64) {
	l.lock(slot).Lock()
}

// Unlock unlocks the mutex of given slot.
func (l *PerSlotLocker) Unlock(slot uint64) {
	l.lock(slot).Unlock()
}
