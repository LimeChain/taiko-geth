package miner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
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

	txSlotSnapshot := miner.worker.txSnapshotsBuilder.getTxSlotSnapshot(slot)
	log.Error("Fetching preconf txs from slot snapshot", "slot", slot, "txs", txSlotSnapshot.Txs)
	txs = append(txs, txSlotSnapshot.Txs...)

	miner.worker.txSnapshotsBuilder.proposeFromTxPoolSnapshot()
	txPoolSnapshot := miner.worker.txSnapshotsBuilder.getTxPoolSnapshot()
	log.Error("Fetching non-preconf txs from txpool snapshot", "txs", txPoolSnapshot.NewTxs)
	txs = append(txs, txPoolSnapshot.NewTxs...)

	// TODO(limechain): support multiple tx lists
	txList := &PreBuiltTxList{
		TxList:           txs,
		EstimatedGasUsed: totalGasUsed,
		BytesLength:      totalBytes,
	}

	return []*PreBuiltTxList{txList}, nil
}
