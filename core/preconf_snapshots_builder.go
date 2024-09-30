package core

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/slocks"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// CHANGE(limechain):

type TxSnapshotsBuilder struct {
	// tx slot snapshot locks
	txSlotSnapshotMu *slocks.PerSlotLocker
	// tx pool snapshot lock
	txPoolSnapshotMu sync.Mutex

	// tx slot snapshots short term cache
	txSlotCache map[uint64]map[common.Hash]bool
	// tx pool snapshots short term cache
	pendingTxCache  map[common.Hash]bool
	proposedTxCache map[common.Hash]bool

	db            ethdb.Database
	preconfTxFeed event.Feed
}

func NewTxSnapshotsBuilder(db ethdb.Database, invPreconfTxEventCh chan InvalidPreconfTxEvent) *TxSnapshotsBuilder {
	b := &TxSnapshotsBuilder{
		txSlotSnapshotMu: &slocks.PerSlotLocker{},
		txSlotCache:      initTxSlotCache(),
		txPoolSnapshotMu: sync.Mutex{},
		pendingTxCache:   initTxCache(),
		proposedTxCache:  initTxCache(),
		db:               db,
	}
	b.preconfTxFeed.Subscribe(invPreconfTxEventCh)
	return b
}

// Preconf txs handling

// EnsureTxSlotSnapshot initializes slot snapshot if it does not exist.
func (b *TxSnapshotsBuilder) EnsureTxSlotSnapshot(slotIndex uint64) {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)
	if txSlotSnapshot == nil {
		rawdb.WriteTxSlotSnapshot(b.db, slotIndex, types.NewTxSlotSnapshot())
	}
}

// GetTxSlotSnapshot retrieves the tx slot snapshot for specific slot.
func (b *TxSnapshotsBuilder) GetTxSlotSnapshot(slotIndex uint64) *types.TxSlotSnapshot {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)
	if txSlotSnapshot == nil {
		return types.NewTxSlotSnapshot()
	}

	return txSlotSnapshot
}

// UpdateTxSlotSnapshot updates the slot snapshot with the new txs.
func (b *TxSnapshotsBuilder) UpdateTxSlotSnapshot(slotIndex uint64, txs []*types.Transaction, bytes uint64, gas uint64) *types.TxSlotSnapshot {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	l1GenesisTimestamp := rawdb.ReadL1GenesisTimestamp(b.db)
	if l1GenesisTimestamp == nil {
		log.Error("Failed to fetch L1 genesis timestamp")
		return nil
	}
	currentSlot, _ := common.CurrentSlotAndEpoch(*l1GenesisTimestamp, time.Now().Unix())

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)

	// TODO(limechain): do not reset past slots until new epoch change
	// since we need previous snapshots for gas/bytes calculations
	// if b.resetPastTxSlotSnapshot(slotIndex, currentSlot) {
	// 	return nil // txSlotSnapshot
	// }

	// Short term cache to speed up the lookup.
	loadTxsInCache(b, slotIndex, txSlotSnapshot)

	for _, tx := range txs {
		// Only preconfirmation txs are stored in the slot snapshot.
		if tx.Type() == types.InclusionPreconfirmationTxType {
			// Remove preconfirmation txs for slot deadlines that have already passed.
			if tx.Deadline().Uint64() < currentSlot {
				b.preconfTxFeed.Send(InvalidPreconfTxEvent{TxHash: tx.Hash()})
				continue
			}

			// Skip inclusion preconfirmation txs that are not for the slot
			// we are currently updating.
			if common.SlotIndex(tx.Deadline().Uint64()) != slotIndex {
				continue
			}

			if !b.txSlotCache[slotIndex][tx.Hash()] {
				txSlotSnapshot.Txs = append(txSlotSnapshot.Txs, tx)
				b.txSlotCache[slotIndex][tx.Hash()] = true
			}
		}
	}

	txSlotSnapshot.GasUsed = gas
	txSlotSnapshot.BytesLength = bytes

	rawdb.WriteTxSlotSnapshot(b.db, slotIndex, txSlotSnapshot)

	return txSlotSnapshot
}

func (b *TxSnapshotsBuilder) UpdateBytesAndGasEstimate(txSlotSnapshot *types.TxSlotSnapshot, bytesUsed uint64, gasUsed uint64) {
	b.txSlotSnapshotMu.Lock(txSlotSnapshot.SlotIndex)
	defer b.txSlotSnapshotMu.Unlock(txSlotSnapshot.SlotIndex)

	txSlotSnapshot.BytesLength += bytesUsed
	txSlotSnapshot.GasUsed += gasUsed
	rawdb.WriteTxSlotSnapshot(b.db, txSlotSnapshot.SlotIndex, txSlotSnapshot)
}

// // resetPastTxSlotSnapshot resets tx snapshot for slot which is in the past.
// func (b *TxSnapshotsBuilder) resetPastTxSlotSnapshot(slotIndex uint64, currentSlot uint64) bool {
// 	if slotIndex < common.SlotIndex(currentSlot) {
// 		b.resetTxSlotSnapshotAndCache(slotIndex)
// 		// log.Warn("Past slot snapshot has been reset", "snapshot slot", slotIndex, "current slot", common.SlotIndex(currentSlot))
// 		return true
// 	}
// 	return false
// }

// ResetTxSlotSnapshot resets tx snapshot for specific slot.
func (b *TxSnapshotsBuilder) ResetTxSlotSnapshot(slotIndex uint64) {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	b.resetTxSlotSnapshotAndCache(slotIndex)
}

// resetTxSlotSnapshotAndCache resets the slot snapshot and cache.
func (b *TxSnapshotsBuilder) resetTxSlotSnapshotAndCache(slotIndex uint64) {
	rawdb.WriteTxSlotSnapshot(b.db, slotIndex, types.NewTxSlotSnapshot())

	if b.txSlotCache != nil {
		b.txSlotCache[slotIndex] = make(map[common.Hash]bool)
	}
}

func loadTxsInCache(b *TxSnapshotsBuilder, slotIndex uint64, slotTxSnapshot *types.TxSlotSnapshot) {
	if (slotTxSnapshot == types.NewTxSlotSnapshot()) || b.txSlotCache == nil {
		b.txSlotCache = initTxSlotCache()
		return
	}

	for _, tx := range slotTxSnapshot.Txs {
		b.txSlotCache[(slotIndex)][tx.Hash()] = true
	}
}

func initTxSlotCache() map[uint64]map[common.Hash]bool {
	cache := make(map[uint64]map[common.Hash]bool)
	for i := 0; i < common.EpochLength; i++ {
		cache[uint64(i)] = make(map[common.Hash]bool)
	}
	return cache
}

// Non-preconf txs handling

// GetTxPoolSnapshot retrieves the tx pool snapshot.
func (b *TxSnapshotsBuilder) GetTxPoolSnapshot() *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)
	if txPoolSnapshot == nil {
		return types.NewTxPoolSnapshot()
	}

	return txPoolSnapshot
}

// UpdateTxPoolSnapshot updates the tx pool snapshot with the new txs.
func (b *TxSnapshotsBuilder) UpdateTxPoolSnapshot(txs []*types.Transaction, bytes uint64, gas uint64) *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)
	if txPoolSnapshot == nil {
		// Initialize the tx pool snapshot
		log.Info("Initialize tx pool snapshot")
		txPoolSnapshot = types.NewTxPoolSnapshot()
	}

	// Short lived cache to speed up the lookup
	loadPendingInCache(b, txPoolSnapshot)

	for _, tx := range txs {
		if tx.Type() != types.InclusionPreconfirmationTxType {
			if !b.pendingTxCache[tx.Hash()] {
				txPoolSnapshot.PendingTxs = append(txPoolSnapshot.PendingTxs, tx)
				b.pendingTxCache[tx.Hash()] = true
			}
		}
	}

	txPoolSnapshot.GasUsed = gas
	txPoolSnapshot.BytesLength = bytes

	rawdb.WriteTxPoolSnapshot(b.db, txPoolSnapshot)

	return txPoolSnapshot
}

// ProposeFromTxPoolSnapshot proposes txs from the tx pool snapshot.
func (b *TxSnapshotsBuilder) ProposeFromTxPoolSnapshot() *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)

	// Short lived cache to speed up the lookup
	loadProposedInCache(b, txPoolSnapshot)

	// Do not reset the 'NewTxs' field to handle the case of 'driver' failures.
	// Each 'propose' call will accumulate until the 'driver' finally executes
	// and resets the state.
	txPoolSnapshot.NewTxs = []*types.Transaction{}

	for _, tx := range txPoolSnapshot.PendingTxs {
		if !b.proposedTxCache[tx.Hash()] {
			txPoolSnapshot.NewTxs = append(txPoolSnapshot.NewTxs, tx)
			txPoolSnapshot.ProposedTxs = append(txPoolSnapshot.ProposedTxs, tx)
			if b.proposedTxCache != nil {
				b.proposedTxCache[tx.Hash()] = true
			} else {
				b.proposedTxCache = map[common.Hash]bool{tx.Hash(): true}
			}
		}
	}

	rawdb.WriteTxPoolSnapshot(b.db, txPoolSnapshot)

	return txPoolSnapshot
}

// RevertProposedTxPoolSnapshot reverts the proposed txs from the tx pool snapshot.
func (b *TxSnapshotsBuilder) RevertProposedTxPoolSnapshot() *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)

	txPoolSnapshot.ProposedTxs = []*types.Transaction{}
	txPoolSnapshot.NewTxs = []*types.Transaction{}

	rawdb.WriteTxPoolSnapshot(b.db, txPoolSnapshot)

	return txPoolSnapshot
}

func (b *TxSnapshotsBuilder) resetTxPoolSnapshot() {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)
	if txPoolSnapshot == nil {
		rawdb.WriteTxPoolSnapshot(b.db, types.NewTxPoolSnapshot())
		return
	}

	newPendingTxs := []*types.Transaction{}

	// Short lived cache to speed up the lookup
	loadProposedInCache(b, txPoolSnapshot)
	for _, tx := range txPoolSnapshot.PendingTxs {
		if !b.proposedTxCache[tx.Hash()] {
			newPendingTxs = append(newPendingTxs, tx)
		}
	}

	// TODO(limechain):
	// Handle the case where the 'driver' has failed and the 'proposer'
	// has continued working, new pending transactions were added after
	// that some txs have been proposed, and the 'driver' has been restarted.
	if len(newPendingTxs) > 0 {
		log.Error("There are new pending txs not in proposed", "txs", newPendingTxs)
	}

	// reset the snapshot and the caches
	rawdb.WriteTxPoolSnapshot(b.db, types.NewTxPoolSnapshot())
	b.pendingTxCache = initTxCache()
	b.proposedTxCache = initTxCache()
}

func loadPendingInCache(b *TxSnapshotsBuilder, txPoolSnapshot *types.TxPoolSnapshot) {
	if txPoolSnapshot == nil || (txPoolSnapshot == types.NewTxPoolSnapshot()) || b.pendingTxCache == nil {
		b.pendingTxCache = initTxCache()
		return
	}

	for _, tx := range txPoolSnapshot.PendingTxs {
		b.pendingTxCache[tx.Hash()] = true
	}
}

func loadProposedInCache(b *TxSnapshotsBuilder, txPoolSnapshot *types.TxPoolSnapshot) {
	if txPoolSnapshot == nil || (txPoolSnapshot == types.NewTxPoolSnapshot()) || b.proposedTxCache == nil {
		b.proposedTxCache = initTxCache()
		return
	}

	for _, tx := range txPoolSnapshot.ProposedTxs {
		b.proposedTxCache[tx.Hash()] = true
	}
}

func initTxCache() map[common.Hash]bool {
	cache := make(map[common.Hash]bool)
	return cache
}

// ResetAllTxSnapshots resets all tx snapshots.
func (b *TxSnapshotsBuilder) ResetAllTxSnapshots() {
	for i := 0; i < common.EpochLength; i++ {
		b.ResetTxSlotSnapshot(uint64(i))
	}
	b.resetTxPoolSnapshot()
	log.Warn("Slot snapshots and txpool snapshot have been reset")
}
