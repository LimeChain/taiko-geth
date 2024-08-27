package rawdb

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	perSlotConstraintsKey = []byte("PerSlotConstraints")
)

func ReadPerSlotConstraints(db ethdb.Database) *types.PerSlotConstraints {
	data, _ := db.Get(perSlotConstraintsKey)
	if len(data) == 0 {
		return nil
	}

	perSlotConstraints := new(types.PerSlotConstraints)
	if err := rlp.DecodeBytes(data, perSlotConstraints); err != nil {
		log.Error("Invalid PerSlotConstraints RLP", "err", err)
		return nil
	}

	return perSlotConstraints
}

func WritePerSlotConstraints(db ethdb.Database, perSlotConstraints *types.PerSlotConstraints) {
	data, err := rlp.EncodeToBytes(perSlotConstraints)
	if err != nil {
		log.Crit("Failed to RLP encode PerSlotConstraints", "err", err)
	}

	if err := db.Put(perSlotConstraintsKey, data); err != nil {
		log.Crit("Failed to store PerSlotConstraints", "err", err)
	}
}
