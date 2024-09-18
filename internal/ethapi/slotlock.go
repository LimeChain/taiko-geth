package ethapi

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

type SlotEstimatesLocker struct {
	mu    sync.Mutex
	locks map[uint64]*sync.Mutex
}

// lock returns the lock of the given slot estimate.
func (l *SlotEstimatesLocker) lock(slot uint64) *sync.Mutex {
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

// LockSlotEstimate locks an slot's mutex.
// This is used to prevent another tx getting and updating the same slot estimates until the lock is released.
func (l *SlotEstimatesLocker) Lock(slot uint64) {
	l.lock(slot).Lock()
}

// UnlockSlotEstimate unlocks the mutex of the given account.
func (l *SlotEstimatesLocker) Unlock(slot uint64) {
	l.lock(slot).Unlock()
}

type StoredSlotEstimates struct {
	db          ethdb.Database
	slotEstLock *SlotEstimatesLocker
}

func NewStoredSlotEstimates(db ethdb.Database, slotEstLock *SlotEstimatesLocker) *StoredSlotEstimates {
	return &StoredSlotEstimates{db, slotEstLock}
}

func (s *StoredSlotEstimates) Read(slot uint64) *types.SlotEstimates {
	s.slotEstLock.Lock(slot)
	defer s.slotEstLock.Unlock(slot)

	slotEstimates := rawdb.ReadSlotEstimates(s.db, slot)
	if slotEstimates == nil {
		slotEstimates = &types.SlotEstimates{}
		rawdb.WriteSlotEstimates(s.db, slot, slotEstimates)
	}
	return slotEstimates
}

func (s *StoredSlotEstimates) Write(slot uint64, slotEstimates *types.SlotEstimates) {
	s.slotEstLock.Lock(slot)
	defer s.slotEstLock.Unlock(slot)

	rawdb.WriteSlotEstimates(s.db, slot, slotEstimates)
}

// TODO(limechain): implement
func (s *StoredSlotEstimates) ResetPast() {}
