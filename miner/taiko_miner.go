package miner

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// PreBuiltTxList is a pre-built transaction list based on the latest chain state,
// with estimated gas used / bytes.
type PreBuiltTxList struct {
	TxList           types.Transactions
	EstimatedGasUsed uint64
	BytesLength      uint64
}

// SealBlockWith mines and seals a block without changing the canonical chain.
func (miner *Miner) SealBlockWith(
	parent common.Hash,
	timestamp uint64,
	blkMeta *engine.BlockMetadata,
	baseFeePerGas *big.Int,
	withdrawals types.Withdrawals,
) (*types.Block, error) {
	return miner.worker.sealBlockWith(parent, timestamp, blkMeta, baseFeePerGas, withdrawals)
}

// CHANGE(limechain):
// FetchTxList retrieves already pre-built list of txs.
func (miner *Miner) FetchTxList(slot uint64) ([]*PreBuiltTxList, error) {
	// TODO(limechain): fetch preconf txs for slot + other txs from the txpool

	log.Warn("Fetching tx list", "slot", slot)
	db := miner.worker.chain.DB()

	l1GenesisTimestamp := rawdb.ReadL1GenesisTimestamp(db)
	if l1GenesisTimestamp == nil {
		return nil, errors.New("failed to fetch L1 genesis timestamp")
	}

	txListConfig := rawdb.ReadTxListConfig(db)
	if txListConfig == nil {
		return nil, errors.New("failed to fetch tx list config")
	}

	var (
		txs          types.Transactions
		totalGasUsed uint64
		totalBytes   uint64
	)

	// currentSlot, _ := common.CurrentSlotAndEpoch(*l1GenesisTimestamp, time.Now().Unix())
	// slotIndex := common.SlotIndex(currentSlot)

	// for i := slotIndex; i < uint64(common.EpochLength); i++ {
	// slotTxSnapshot := rawdb.ReadSlotTxSnapshot(db, i)
	// if totalGasUsed+slotTxSnapshot.GasUsed > txListConfig.BlockMaxGasLimit ||
	// 	totalBytes+slotTxSnapshot.BytesLength > txListConfig.MaxBytesPerTxList {
	// 	break
	// }

	slotTxSnapshot := miner.worker.ProposeSlotSnapshotTxs(slot)
	// if len(slotTxSnapshot.NewTxs) == 0 {
	// break
	// }

	txs = append(txs, slotTxSnapshot.NewTxs...)
	totalGasUsed += slotTxSnapshot.GasUsed
	totalBytes += slotTxSnapshot.BytesLength
	// }

	// TODO(limechain): support multiple tx lists
	txList := &PreBuiltTxList{
		TxList:           txs,
		EstimatedGasUsed: totalGasUsed,
		BytesLength:      totalBytes,
	}

	return []*PreBuiltTxList{txList}, nil
}
