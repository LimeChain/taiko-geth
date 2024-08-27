package rawdb

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	assignedL1Slots = []byte("AssignedL1Slots")
)

func ReadAssignedL1Slots(db ethdb.Database) []uint64 {
	data, err := db.Get(assignedL1Slots)
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	assignedSlots := make([]uint64, 0)
	if err := rlp.DecodeBytes(data, &assignedSlots); err != nil {
		log.Error("Invalid AssignedL1Slots RLP", "err", err)
		return nil
	}

	return assignedSlots
}

func WriteAssignedL1Slots(db ethdb.Database, slots []uint64) {
	data, err := rlp.EncodeToBytes(slots)
	if err != nil {
		log.Crit("Failed to RLP encode AssignedL1Slots", "err", err)
	}

	if err := db.Put(assignedL1Slots, data); err != nil {
		log.Crit("Failed to store AssignedL1Slots", "err", err)
	}
}
