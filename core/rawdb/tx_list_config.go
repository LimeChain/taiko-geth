package rawdb

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	txListConfigKey = []byte("TxListConfig")
)

// TODO(limechain): move in types
type TxListConfig struct {
	BaseFee           *big.Int // calculated in the protocol contract
	BlockMaxGasLimit  uint64   // hard-coded in the protocol contract config
	MaxBytesPerTxList uint64
}

func ReadTxListConfig(db ethdb.Database) *TxListConfig {
	data, _ := db.Get(txListConfigKey)
	if len(data) == 0 {
		return nil
	}

	txListConfig := new(TxListConfig)
	if err := rlp.DecodeBytes(data, txListConfig); err != nil {
		log.Error("Invalid TxListConfig RLP", "err", err)
		return nil
	}

	return txListConfig
}

func WriteTxListConfig(db ethdb.Database, txListConfig *TxListConfig) {
	data, err := rlp.EncodeToBytes(txListConfig)
	if err != nil {
		log.Crit("Failed to RLP encode TxListConfig", "err", err)
	}

	if err := db.Put(txListConfigKey, data); err != nil {
		log.Crit("Failed to store TxListConfig", "err", err)
	}
}
