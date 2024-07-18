package miner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

// BuildTransactionsLists builds multiple transactions lists which satisfy all the given limits.
func (miner *Miner) BuildTransactionList(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) error {
	err := miner.worker.BuildTransactionList(
		beneficiary,
		baseFee,
		blockMaxGasLimit,
		maxBytesPerTxList,
		locals,
		maxTransactionsLists,
	)
	return err
}

func (miner *Miner) FetchTransactionList(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) ([]*types.PreBuiltTxList, error) {
	// TODO: handle multiple tx lists
	txListState := miner.worker.ProposeFromTxListState()

	txList := &types.PreBuiltTxList{
		TxList:           txListState.NewTxs,
		EstimatedGasUsed: txListState.EstimatedGasUsed,
		BytesLength:      txListState.BytesLength,
	}

	return []*types.PreBuiltTxList{txList}, nil
}
