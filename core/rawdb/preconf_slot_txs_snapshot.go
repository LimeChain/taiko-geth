package rawdb

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	slotTxSnapshotPrefixKey = []byte("SlotTxSnapshot")
)

func ReadSlotTxSnapshot(db ethdb.Database, slot uint64) *types.SlotTxSnapshot {
	data, _ := db.Get(slotTxSnapshotKey(slot))
	if len(data) == 0 {
		return nil
	}

	slotTxSnapshot := new(types.SlotTxSnapshot)
	if err := rlp.Decode(bytes.NewBuffer(data), slotTxSnapshot); err != nil {
		log.Error("Invalid SlotTxSnapshot RLP", "err", err)
		return nil
	}

	return slotTxSnapshot
}

func WriteSlotTxSnapshot(db ethdb.Database, slot uint64, slotTxSnapshot *types.SlotTxSnapshot) {
	data := bytes.NewBuffer(nil)
	err := rlp.Encode(data, slotTxSnapshot)
	if err != nil {
		log.Crit("Failed to RLP encode SlotTxSnapshot", "err", err)
	}

	if err := db.Put(slotTxSnapshotKey(slot), data.Bytes()); err != nil {
		log.Crit("Failed to store SlotTxSnapshot", "err", err)
	}
}

func slotTxSnapshotKey(slot uint64) []byte {
	index := slot % uint64(common.EpochLength)
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, index)
	return append(slotTxSnapshotPrefixKey, enc...)
}
