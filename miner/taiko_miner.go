package miner

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// FetchTransactionList retrieves already pre-built list of txs.
func (miner *Miner) FetchTransactionList() ([]*PreBuiltTxList, error) {
	txPoolSnapshot, perSlotConstraints := miner.worker.ProposeTxsInPoolSnapshot()

	if txPoolSnapshot == nil || perSlotConstraints == nil {
		return nil, errors.New("failed to fetch txs to propose")
	}

	// TODO(limechain): refactor, no need to return multiple lists
	txList := &PreBuiltTxList{
		TxList:           txPoolSnapshot.NewTxs,
		EstimatedGasUsed: perSlotConstraints.Total.EstimatedGasUsed,
		BytesLength:      perSlotConstraints.Total.BytesLength,
	}

	return []*PreBuiltTxList{txList}, nil
}
