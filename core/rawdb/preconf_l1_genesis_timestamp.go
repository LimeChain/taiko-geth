package rawdb

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// CHANGE(limechain):

var (
	l1GenesisTimestampKey = []byte("L1GenesisTimestamp")
)

func ReadL1GenesisTimestamp(db ethdb.Database) *uint64 {
	data, err := db.Get(l1GenesisTimestampKey)
	if err != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}

	l1GenesisTimestamp := new(uint64)
	if err := rlp.DecodeBytes(data, l1GenesisTimestamp); err != nil {
		log.Error("Invalid L1GenesisTimestamp RLP", "err", err)
		return nil
	}

	return l1GenesisTimestamp
}

func WriteL1GenesisTimestamp(db ethdb.Database, timestamp uint64) {
	data, err := rlp.EncodeToBytes(timestamp)
	if err != nil {
		log.Crit("Failed to RLP encode L1GenesisTimestamp", "err", err)
	}

	if err := db.Put(l1GenesisTimestampKey, data); err != nil {
		log.Crit("Failed to store L1GenesisTimestamp", "err", err)
	}
}
