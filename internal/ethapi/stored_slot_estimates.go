package ethapi

import (
	"github.com/ethereum/go-ethereum/common/slocks"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

// CHANGE(limechain): Provides a way to store and retrieve gas/bytes
// estimates per slot used during preconf tx submission.

type StoredSlotEstimates struct {
	db       ethdb.Database
	slotLock *slocks.PerSlotLocker
}

func NewStoredSlotEstimates(db ethdb.Database, slotLock *slocks.PerSlotLocker) *StoredSlotEstimates {
	return &StoredSlotEstimates{db, slotLock}
}

func (s *StoredSlotEstimates) Read(slot uint64) *types.SlotEstimates {
	s.slotLock.Lock(slot)
	defer s.slotLock.Unlock(slot)

	slotEstimates := rawdb.ReadSlotEstimates(s.db, slot)
	if slotEstimates == nil {
		slotEstimates = &types.SlotEstimates{}
		rawdb.WriteSlotEstimates(s.db, slot, slotEstimates)
	}
	return slotEstimates
}

func (s *StoredSlotEstimates) Write(slot uint64, slotEstimates *types.SlotEstimates) {
	s.slotLock.Lock(slot)
	defer s.slotLock.Unlock(slot)

	rawdb.WriteSlotEstimates(s.db, slot, slotEstimates)
}
