package eth

import (
	"bytes"
	"compress/zlib"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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

// TxPoolContent retrieves the transaction pool content with the given upper limits.
func (a *TaikoAuthAPIBackend) TxPoolContent(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	locals []string,
	maxTransactionsLists uint64,
) ([]*miner.PreBuiltTxList, error) {
	log.Debug(
		"Fetching L2 pending transactions finished",
		"baseFee", baseFee,
		"blockMaxGasLimit", blockMaxGasLimit,
		"maxBytesPerTxList", maxBytesPerTxList,
		"maxTransactions", maxTransactionsLists,
		"locals", locals,
	)

	return a.eth.Miner().BuildTransactionsLists(
		beneficiary,
		baseFee,
		blockMaxGasLimit,
		maxBytesPerTxList,
		locals,
		maxTransactionsLists,
	)
}

func (a *TaikoAuthAPIBackend) PreconfirmedTxs() ([]*miner.PreBuiltTxList, error) {
	log.Debug("Fetching L2 preconfirmed txs")

	blockHash, blockNum := rawdb.ReadPendingVirtualBlock(a.eth.ChainDb())
	if blockNum == nil {
		return []*miner.PreBuiltTxList{}, nil
	}

	block := rawdb.ReadBlock(a.eth.ChainDb(), blockHash, *blockNum)
	if block == nil {
		log.Debug("Empty virtual block", "blockHash", blockHash, "blockNum", blockNum)
		return []*miner.PreBuiltTxList{}, nil
	}

	b, err := encodeAndComporeessTxList(block.Transactions())
	if err != nil {
		log.Error("Failed to compress block txs", "blockHash", blockHash, "blockNum", blockNum, "err", err)
		return nil, err
	}

	prebuildTxList := &miner.PreBuiltTxList{
		TxList:           block.Transactions(),
		EstimatedGasUsed: block.GasUsed(),
		BytesLength:      uint64(len(b)),
	}

	return []*miner.PreBuiltTxList{prebuildTxList}, nil
}

func (a *TaikoAuthAPIBackend) ProposePreconfirmedTxs() ([]*miner.PreBuiltTxList, error) {
	log.Debug("Fetching L2 preconfirmed txs for proposing")

	blockHash, blockNum := rawdb.ReadPendingVirtualBlock(a.eth.ChainDb())
	if blockNum == nil {
		return []*miner.PreBuiltTxList{}, nil
	}

	block := rawdb.ReadBlock(a.eth.ChainDb(), blockHash, *blockNum)
	if block == nil {
		log.Debug("Empty virtual block", "blockHash", blockHash, "blockNum", blockNum)
		return []*miner.PreBuiltTxList{}, nil
	}

	b, err := encodeAndComporeessTxList(block.Transactions())
	if err != nil {
		log.Error("Failed to compress block txs", "blockHash", blockHash, "blockNum", blockNum, "err", err)
		return nil, err
	}

	prebuildTxList := &miner.PreBuiltTxList{
		TxList:           block.Transactions(),
		EstimatedGasUsed: block.GasUsed(),
		BytesLength:      uint64(len(b)),
	}

	return []*miner.PreBuiltTxList{prebuildTxList}, nil
}

// encodeAndComporeessTxList encodes and compresses the given transactions list.
func encodeAndComporeessTxList(txs types.Transactions) ([]byte, error) {
	b, err := rlp.EncodeToBytes(txs)
	if err != nil {
		return nil, err
	}

	return compress(b)
}

// compress compresses the given txList bytes using zlib.
func compress(txListBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	if _, err := w.Write(txListBytes); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

type HashAndNumber struct {
	Hash   common.Hash
	Number uint64
}

// GetPreconfirmedVirtualBlock ...
func (s *TaikoAPIBackend) GetPreconfirmedVirtualBlock() HashAndNumber {
	hash, number := rawdb.ReadPreconfirmedVirtualBlock(s.eth.ChainDb())

	if (hash == common.Hash{}) && number == nil {
		return HashAndNumber{}
	}

	return HashAndNumber{hash, *number}
}

// GetPendingVirtualBlock ...
func (s *TaikoAPIBackend) GetPendingVirtualBlock() HashAndNumber {
	hash, number := rawdb.ReadPendingVirtualBlock(s.eth.ChainDb())

	if (hash == common.Hash{}) && number == nil {
		return HashAndNumber{}
	}

	return HashAndNumber{hash, *number}
}

// UpdatePreconfirmedVirtualBlock ...
func (s *TaikoAPIBackend) UpdatePreconfirmedVirtualBlock(blockHash common.Hash, blockNumber *math.HexOrDecimal256) bool {
	rawdb.WritePreconfirmedVirtualBlock(s.eth.ChainDb(), blockHash, (*big.Int)(blockNumber))
	return true
}

// UpdatePendingVirtualBlock ...
func (s *TaikoAPIBackend) UpdatePendingVirtualBlock(blockHash common.Hash, blockNumber *math.HexOrDecimal256) bool {
	rawdb.WritePendingVirtualBlock(s.eth.ChainDb(), blockHash, (*big.Int)(blockNumber))

	err := s.eth.blockchain.SetPreconfirmedBlock()
	if err != nil {
		log.Error("failed to set preconfirmed pending block", "error", err)
	}

	return true
}

func (s *TaikoAPIBackend) DeletePendingVirtualBlock() {
	rawdb.DeletePendingVirtualBlock(s.eth.ChainDb())
}
