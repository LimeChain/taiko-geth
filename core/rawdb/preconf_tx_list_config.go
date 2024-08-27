package rawdb

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	txListConfigKey = []byte("TxListConfig")
)

func ReadTxListConfig(db ethdb.Database) *types.TxListConfig {
	data, _ := db.Get(txListConfigKey)
	if len(data) == 0 {
		return nil
	}

	txListConfig := new(types.TxListConfig)
	if err := rlp.DecodeBytes(data, txListConfig); err != nil {
		log.Error("Invalid TxListConfig RLP", "err", err)
		return nil
	}

	return txListConfig
}

func WriteTxListConfig(db ethdb.Database, txListConfig *types.TxListConfig) {
	data, err := rlp.EncodeToBytes(txListConfig)
	if err != nil {
		log.Crit("Failed to RLP encode TxListConfig", "err", err)
	}

	if err := db.Put(txListConfigKey, data); err != nil {
		log.Crit("Failed to store TxListConfig", "err", err)
	}
}
