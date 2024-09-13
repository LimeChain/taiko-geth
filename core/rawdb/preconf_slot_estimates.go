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
	slotEstimatesPrefixKey = []byte("SlotEstimates")
)

func ReadSlotEstimates(db ethdb.Database, slot uint64) *types.SlotEstimates {
	data, _ := db.Get(slotEstimatesKey(slot))
	if len(data) == 0 {
		return nil
	}

	slotEstimates := new(types.SlotEstimates)
	if err := rlp.Decode(bytes.NewBuffer(data), slotEstimates); err != nil {
		log.Error("Invalid SlotEstimates RLP", "err", err)
		return nil
	}

	return slotEstimates
}

func WriteSlotEstimates(db ethdb.Database, slot uint64, slotEstimates *types.SlotEstimates) {
	data := bytes.NewBuffer(nil)
	err := rlp.Encode(data, slotEstimates)
	if err != nil {
		log.Crit("Failed to RLP encode SlotEstimates", "err", err)
	}

	if err := db.Put(slotEstimatesKey(slot), data.Bytes()); err != nil {
		log.Crit("Failed to store SlotEstimates", "err", err)
	}
}

func slotEstimatesKey(slot uint64) []byte {
	index := slot % uint64(common.EpochLength)
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, index)
	return append(slotEstimatesPrefixKey, enc...)
}
