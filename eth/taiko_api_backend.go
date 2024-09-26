package eth

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
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

// CHANGE(limechain):

func (a *TaikoAPIBackend) FetchL1GenesisTimestamp() uint64 {
	timestamp := rawdb.ReadL1GenesisTimestamp(a.eth.ChainDb())
	if timestamp == nil {
		return 0
	}
	return *timestamp
}

func (a *TaikoAPIBackend) FetchAssignedSlots() []uint64 {
	return rawdb.ReadAssignedL1Slots(a.eth.ChainDb())
}

// UpdateConfigAndSlots updates the assigned slots and configuration for
// preaparing tx lists.
func (a *TaikoAuthAPIBackend) UpdateConfigAndSlots(
	l1GenesisTimestamp uint64,
	newAssignedSlots []uint64,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	beneficiary common.Address,
	locals []string,
	maxTransactionsLists uint64,
) error {
	db := a.eth.ChainDb()

	txListConfig := &types.TxListConfig{
		Beneficiary:          beneficiary,
		BaseFee:              baseFee,
		BlockMaxGasLimit:     blockMaxGasLimit,
		MaxBytesPerTxList:    maxBytesPerTxList,
		Locals:               locals,
		MaxTransactionsLists: maxTransactionsLists,
	}

	rawdb.WriteL1GenesisTimestamp(db, l1GenesisTimestamp)
	rawdb.WriteTxListConfig(db, txListConfig)

	currentSlot, _ := common.CurrentSlotAndEpoch(l1GenesisTimestamp, time.Now().Unix())

	storedSlots := rawdb.ReadAssignedL1Slots(db)
	slotsCache := make(map[uint64]bool)
	updatedAssignedSlots := make([]uint64, 0)

	for _, slot := range storedSlots {
		if slot < currentSlot {
			continue
		}

		if _, ok := slotsCache[slot]; !ok {
			updatedAssignedSlots = append(updatedAssignedSlots, slot)
			slotsCache[slot] = true
		}
	}

	for _, slot := range newAssignedSlots {
		if slot < currentSlot {
			continue
		}

		if _, ok := slotsCache[slot]; !ok {
			updatedAssignedSlots = append(updatedAssignedSlots, slot)
			slotsCache[slot] = true
		}
	}

	rawdb.WriteAssignedL1Slots(db, updatedAssignedSlots)
	log.Error("Current assigned slots", "slots", updatedAssignedSlots)

	return nil
}

// FetchTxList retrieves already prepared list of txs.
func (a *TaikoAuthAPIBackend) FetchTxList(slot uint64) ([]*miner.PreBuiltTxList, error) {
	return a.eth.Miner().FetchTxList(slot)
}
