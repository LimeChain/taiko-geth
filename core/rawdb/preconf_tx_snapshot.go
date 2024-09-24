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
	txSlotSnapshotPrefixKey = []byte("TxSlotSnapshot")
	txPoolSnapshotPrefixKey = []byte("TxPoolSnapshot")
)

func ReadTxSlotSnapshot(db ethdb.Database, slot uint64) *types.TxSlotSnapshot {
	data, _ := db.Get(txSlotSnapshotKey(slot))
	if len(data) == 0 {
		return nil
	}

	txSlotSnapshot := new(types.TxSlotSnapshot)
	if err := rlp.Decode(bytes.NewBuffer(data), txSlotSnapshot); err != nil {
		log.Error("Invalid TxSlotSnapshot RLP", "err", err)
		return nil
	}

	return txSlotSnapshot
}

func WriteTxSlotSnapshot(db ethdb.Database, slot uint64, txSlotSnapshot *types.TxSlotSnapshot) {
	data := bytes.NewBuffer(nil)
	err := rlp.Encode(data, txSlotSnapshot)
	if err != nil {
		log.Crit("Failed to RLP encode TxSlotSnapshot", "err", err)
	}

	if err := db.Put(txSlotSnapshotKey(slot), data.Bytes()); err != nil {
		log.Crit("Failed to store TxSlotSnapshot", "err", err)
	}
}

func txSlotSnapshotKey(slot uint64) []byte {
	index := slot % uint64(common.EpochLength)
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, index)
	return append(txSlotSnapshotPrefixKey, enc...)
}

func ReadTxPoolSnapshot(db ethdb.Database) *types.TxPoolSnapshot {
	data, _ := db.Get(txPoolSnapshotPrefixKey)
	if len(data) == 0 {
		return nil
	}

	txPoolSnapshot := new(types.TxPoolSnapshot)
	if err := rlp.Decode(bytes.NewBuffer(data), txPoolSnapshot); err != nil {
		log.Error("Invalid TxPoolSnapshot RLP", "err", err)
		return nil
	}

	return txPoolSnapshot
}

func WriteTxPoolSnapshot(db ethdb.Database, txPoolSnapshot *types.TxPoolSnapshot) {
	data := bytes.NewBuffer(nil)
	err := rlp.Encode(data, txPoolSnapshot)
	if err != nil {
		log.Crit("Failed to RLP encode TxPoolSnapshot", "err", err)
	}

	if err := db.Put(txPoolSnapshotPrefixKey, data.Bytes()); err != nil {
		log.Crit("Failed to store TxPoolSnapshot", "err", err)
	}
}
