package rawdb

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func ReadTxListState(db ethdb.Database) *types.TxListState {
	data, _ := db.Get(txListStateKey(0))
	if len(data) == 0 {
		return nil
	}

	txListState := new(types.TxListState)
	if err := rlp.DecodeBytes(data, txListState); err != nil {
		log.Error("Invalid TxListState RLP", "err", err)
		return nil
	}

	return txListState
}

func WriteTxListState(db ethdb.Database, txListState *types.TxListState) {
	data, err := rlp.EncodeToBytes(txListState)
	if err != nil {
		log.Crit("Failed to RLP encode TxListState", "err", err)
	}

	if err := db.Put(txListStateKey(0), data); err != nil {
		log.Crit("Failed to store TxListState", "err", err)
	}
}
