package rawdb

import (
	"bytes"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	l1EpochKey         = []byte("l1CurrentEpoch")
	l1SlotKey          = []byte("l1CurrentSlot")
	l1AssignedSlotsKey = []byte("l1AssignedSlots")
)

func ReadCurrentL1Epoch(db ethdb.KeyValueReader) uint64 {
	data, err := db.Get(l1EpochKey)
	if err != nil {
		return 0
	}
	if len(data) == 0 {
		return 0
	}

	var epoch uint64
	err = rlp.DecodeBytes(data, &epoch)
	if err != nil {
		log.Crit("Failed to decode epoch", "err", err)
	}

	return epoch
}

func WriteCurrentL1Epoch(db ethdb.KeyValueWriter, epoch uint64) {
	data := bytes.NewBuffer(nil)

	err := rlp.Encode(data, epoch)
	if err != nil {
		log.Crit("Failed to RLP encode epoch", "err", err)
	}

	if err := db.Put(l1EpochKey, data.Bytes()); err != nil {
		log.Crit("Failed to store L1 epoch", "error", err)
	}
}

func ReadCurrentL1Slot(db ethdb.KeyValueReader) uint64 {
	data, err := db.Get(l1SlotKey)
	if err != nil {
		return 0
	}
	if len(data) == 0 {
		return 0
	}

	var slot uint64
	err = rlp.DecodeBytes(data, &slot)
	if err != nil {
		log.Crit("Failed to decode slot", "err", err)
	}

	return slot
}

func WriteCurrentL1Slot(db ethdb.KeyValueWriter, slot uint64) {
	data := bytes.NewBuffer(nil)

	err := rlp.Encode(data, slot)
	if err != nil {
		log.Crit("Failed to RLP encode slot", "err", err)
	}

	if err := db.Put(l1SlotKey, data.Bytes()); err != nil {
		log.Crit("Failed to store L1 slot", "error", err)
	}
}

func ReadAssignedL1Slots(db ethdb.Database) []uint64 {
	data, err := db.Get(l1AssignedSlotsKey)
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	assignedSlots := make([]uint64, 0)
	err = rlp.DecodeBytes(data, &assignedSlots)
	if err != nil {
		log.Crit("Failed to decode assigned slots", "err", err)
	}

	return assignedSlots
}

func WriteAssignedL1Slots(db ethdb.Database, slots []uint64) {
	data := bytes.NewBuffer(nil)

	err := rlp.Encode(data, slots)
	if err != nil {
		log.Crit("Failed to RLP encode assigned slots", "err", err)
	}

	if err := db.Put(l1AssignedSlotsKey, data.Bytes()); err != nil {
		log.Crit("Failed to store assigned slots", "err", err)
	}
}
