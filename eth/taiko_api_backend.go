package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
)

// TaikoAPIBackend handles L2 node related RPC calls.
type TaikoAPIBackend struct {
	eth *Ethereum
}

// NewTaikoAPIBackend creates a new TaikoAPIBackend instance.
func NewTaikoAPIBackend(eth *Ethereum) *TaikoAPIBackend {
	return &TaikoAPIBackend{
		eth: eth,
	}
}

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (s *TaikoAPIBackend) HeadL1Origin() (*rawdb.L1Origin, error) {
	blockID, err := rawdb.ReadHeadL1Origin(s.eth.ChainDb())
	if err != nil {
		return nil, err
	}

	if blockID == nil {
		return nil, ethereum.NotFound
	}

	l1Origin, err := rawdb.ReadL1Origin(s.eth.ChainDb(), blockID)
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}

// L1OriginByID returns the L2 block's corresponding L1 origin.
func (s *TaikoAPIBackend) L1OriginByID(blockID *math.HexOrDecimal256) (*rawdb.L1Origin, error) {
	l1Origin, err := rawdb.ReadL1Origin(s.eth.ChainDb(), (*big.Int)(blockID))
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}

// GetSyncMode returns the node sync mode.
func (s *TaikoAPIBackend) GetSyncMode() (string, error) {
	return s.eth.config.SyncMode.String(), nil
}

// TaikoAuthAPIBackend handles L2 node related authorized RPC calls.
type TaikoAuthAPIBackend struct {
	eth *Ethereum
}

// NewTaikoAuthAPIBackend creates a new TaikoAuthAPIBackend instance.
func NewTaikoAuthAPIBackend(eth *Ethereum) *TaikoAuthAPIBackend {
	return &TaikoAuthAPIBackend{eth}
}

// BuildTxList initiates the process of building tx lists.
func (a *TaikoAuthAPIBackend) BuildTxList(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) error {
	err := a.eth.Miner().BuildTransactionList(
		beneficiary,
		baseFee,
		blockMaxGasLimit,
		maxBytesPerTxList,
		locals,
		maxTransactionsLists,
	)
	return err
}

// FetchTxList retrieves already pre-built list of txs.
func (a *TaikoAuthAPIBackend) FetchTxList(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) ([]*miner.PreBuiltTxList, error) {
	log.Info(
		"Fetching L2 transactions to propose",
		"baseFee", baseFee,
		"blockMaxGasLimit", blockMaxGasLimit,
		"maxBytesPerTxList", maxBytesPerTxList,
		"maxTransactions", maxTransactionsLists,
		"locals", locals,
	)
	return a.eth.Miner().FetchTransactionList(
		beneficiary,
		baseFee,
		blockMaxGasLimit,
		maxBytesPerTxList,
		locals,
		maxTransactionsLists,
	)
}
