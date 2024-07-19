package miner

import (
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

// BuildTransactionsLists initiates the process of building tx lists that later can be fetched.
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

// FetchTransactionList retrieves already pre-built list of txs.
func (miner *Miner) FetchTransactionList(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) ([]*PreBuiltTxList, error) {
	txPoolSnapshot := miner.worker.ProposeTxsInPoolSnapshot()

	// TODO(limechain): handle multiple tx lists

	txList := &PreBuiltTxList{
		TxList:           txPoolSnapshot.NewTxs,
		EstimatedGasUsed: txPoolSnapshot.EstimatedGasUsed,
		BytesLength:      txPoolSnapshot.BytesLength,
	}

	return []*PreBuiltTxList{txList}, nil
}
