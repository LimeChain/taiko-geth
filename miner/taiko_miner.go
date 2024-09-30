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

// FetchTxList retrieves already pre-built list of preconf txs for specific slot and
// non-preconf txs currently in the txpool.
func (miner *Miner) FetchTxList(slot uint64) ([]*PreBuiltTxList, error) {
	var (
		txs          types.Transactions
		totalGasUsed uint64
		totalBytes   uint64
	)

	db := miner.worker.chain.DB()

	txListConfig := rawdb.ReadTxListConfig(db)
	if txListConfig == nil {
		return nil, errors.New("failed to fetch tx list config")
	}

	slotIndex := common.SlotIndex(slot)
	txSlotSnapshot := miner.worker.txSnapshotsBuilder.GetTxSlotSnapshot(slotIndex)

	// Fetch preconf txs from slot snapshot
	if txSlotSnapshot.GasUsed < txListConfig.BlockMaxGasLimit &&
		txSlotSnapshot.BytesLength < txListConfig.MaxBytesPerTxList {

		txs = append(txs, txSlotSnapshot.Txs...)
		totalGasUsed += txSlotSnapshot.GasUsed
		totalBytes += txSlotSnapshot.BytesLength
	} else {
		log.Error("Tx list limits reached", "slot index", common.SlotIndex(slot), "txs", txs, "gas used", totalGasUsed, "bytes length", totalBytes)
		return nil, nil
	}

	// Fetch non-preconf txs from txpool snapshot
	miner.worker.txSnapshotsBuilder.ProposeFromTxPoolSnapshot()
	txPoolSnapshot := miner.worker.txSnapshotsBuilder.GetTxPoolSnapshot()
	if totalGasUsed+txPoolSnapshot.GasUsed < txListConfig.BlockMaxGasLimit &&
		totalBytes+txPoolSnapshot.BytesLength < txListConfig.MaxBytesPerTxList {

		totalGasUsed += txPoolSnapshot.GasUsed
		totalBytes += txPoolSnapshot.BytesLength
		txs = append(txs, txPoolSnapshot.NewTxs...)
	} else {
		miner.worker.txSnapshotsBuilder.RevertProposedTxPoolSnapshot()
		log.Error("Tx list limits reached", "slot index", common.SlotIndex(slot), "txs", txs, "gas used", totalGasUsed, "bytes length", totalBytes)
	}

	// TODO(limechain): support multiple tx lists
	txList := &PreBuiltTxList{
		TxList:           txs,
		EstimatedGasUsed: totalGasUsed,
		BytesLength:      totalBytes,
	}

	log.Error("Fetch tx list", "slot index", common.SlotIndex(slot), "txs", txs, "gas used", totalGasUsed, "bytes length", totalBytes)
	return []*PreBuiltTxList{txList}, nil
}
