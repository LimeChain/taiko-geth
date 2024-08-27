package rawdb

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	txPoolSnapshotPrefixKey = []byte("TxPoolSnapshot")
)

func ReadTxPoolSnapshot(db ethdb.Database) *types.TxPoolSnapshot {
	data, _ := db.Get(txPoolSnapshotPrefixKey)
	if len(data) == 0 {
		return nil
	}

	txPoolSnapshot := new(types.TxPoolSnapshot)
	if err := rlp.DecodeBytes(data, txPoolSnapshot); err != nil {
		log.Error("Invalid TxPoolSnapshot RLP", "err", err)
		return nil
	}

	return txPoolSnapshot
}

func WriteTxPoolSnapshot(db ethdb.Database, txPoolSnapshot *types.TxPoolSnapshot) {
	data, err := rlp.EncodeToBytes(txPoolSnapshot)
	if err != nil {
		log.Crit("Failed to RLP encode TxPoolSnapshot", "err", err)
	}

	if err := db.Put(txPoolSnapshotPrefixKey, data); err != nil {
		log.Crit("Failed to store TxPoolSnapshot", "err", err)
	}
}
