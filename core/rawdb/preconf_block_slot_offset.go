package rawdb

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	L2BlockToL1SlotOffsetKey = []byte("L2BlockToL1SlotOffset")
)

func ReadL2BlockToL1SlotOffset(db ethdb.Database) *uint64 {
	data, err := db.Get(L2BlockToL1SlotOffsetKey)
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	offset := new(uint64)
	if err := rlp.DecodeBytes(data, offset); err != nil {
		log.Error("Invalid L2BlockToL1SlotOffset RLP", "err", err)
		return nil
	}

	return offset
}

func WriteL2BlockToL1SlotOffset(db ethdb.Database, offset uint64) {
	data, err := rlp.EncodeToBytes(offset)
	if err != nil {
		log.Crit("Failed to RLP encode L2BlockToL1SlotOffset", "err", err)
	}

	if err := db.Put(L2BlockToL1SlotOffsetKey, data); err != nil {
		log.Crit("Failed to store L2BlockToL1SlotOffset", "err", err)
	}
}
