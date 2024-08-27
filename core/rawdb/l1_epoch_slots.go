package rawdb

import (
	"bytes"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	l1AssignedSlotsKey = []byte("assignedL1Slots")
)

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
