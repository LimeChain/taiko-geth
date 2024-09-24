package miner

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/slocks"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

// CHANGE(limechain):

type txSnapshotsBuilder struct {
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

func newTxSnapshotsBuilder(db ethdb.Database, invPreconfTxEventCh chan core.InvalidPreconfTxEvent) *txSnapshotsBuilder {
	b := &txSnapshotsBuilder{
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

// initTxSlotSnapshot initializes slot snapshot if it does not exist.
func (b *txSnapshotsBuilder) initTxSlotSnapshot(slotIndex uint64) {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)
	if txSlotSnapshot == nil {
		rawdb.WriteTxSlotSnapshot(b.db, slotIndex, types.NewTxSlotSnapshot())
	}
}

func (b *txSnapshotsBuilder) getTxSlotSnapshot(slotIndex uint64) *types.TxSlotSnapshot {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)
	if txSlotSnapshot == nil {
		return types.NewTxSlotSnapshot()
	}

	return txSlotSnapshot
}

// updateSlotSnapshotTxs updates the slot snapshot with the new txs.
func (b *txSnapshotsBuilder) updateTxSlotSnapshot(slotIndex uint64, txs []*types.Transaction, bz []byte, env *environment) *types.TxSlotSnapshot {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	l1GenesisTimestamp := rawdb.ReadL1GenesisTimestamp(b.db)
	if l1GenesisTimestamp == nil {
		log.Error("Failed to fetch L1 genesis timestamp")
		return nil
	}
	currentSlot, _ := common.CurrentSlotAndEpoch(*l1GenesisTimestamp, time.Now().Unix())

	txSlotSnapshot := rawdb.ReadTxSlotSnapshot(b.db, slotIndex)

	if b.resetPastTxSlotSnapshot(slotIndex, currentSlot) {
		// TODO(limechain):
		return nil // txSlotSnapshot
	}

	// Short term cache to speed up the lookup.
	loadTxsInCache(b, slotIndex, txSlotSnapshot)

	for _, tx := range txs {
		// Only preconfirmation txs are stored in the slot snapshot.
		if tx.Type() == types.InclusionPreconfirmationTxType {
			// Remove preconfirmation txs for slot deadlines that have already passed.
			if tx.Deadline().Uint64() < currentSlot {
				b.preconfTxFeed.Send(core.InvalidPreconfTxEvent{TxHash: tx.Hash()})
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

	// TODO(limechain): fix calculation to be for preconf txs
	txSlotSnapshot.GasUsed = env.header.GasLimit - env.gasPool.Gas()
	txSlotSnapshot.BytesLength = uint64(len(bz))
	rawdb.WriteTxSlotSnapshot(b.db, slotIndex, txSlotSnapshot)

	return txSlotSnapshot
}

// resetTxSlotSnapshot resets tx snapshot for specific slot.
func (b *txSnapshotsBuilder) resetTxSlotSnapshot(slotIndex uint64) {
	b.txSlotSnapshotMu.Lock(slotIndex)
	defer b.txSlotSnapshotMu.Unlock(slotIndex)

	b.resetTxSlotSnapshotAndCache(slotIndex)
}

// resetPastTxSlotSnapshot resets tx snapshot for slot which is in the past.
func (b *txSnapshotsBuilder) resetPastTxSlotSnapshot(slotIndex uint64, currentSlot uint64) bool {
	if slotIndex < common.SlotIndex(currentSlot) {
		b.resetTxSlotSnapshotAndCache(slotIndex)
		log.Warn("Past slot snapshot has been reset", "snapshot slot", slotIndex, "current slot", currentSlot)
		return true
	}
	return false
}

// resetTxSlotSnapshotAndCache resets the slot snapshot and cache.
func (b *txSnapshotsBuilder) resetTxSlotSnapshotAndCache(slotIndex uint64) {
	rawdb.WriteTxSlotSnapshot(b.db, slotIndex, types.NewTxSlotSnapshot())

	if b.txSlotCache != nil {
		b.txSlotCache[slotIndex] = make(map[common.Hash]bool)
	}
}

func loadTxsInCache(b *txSnapshotsBuilder, slotIndex uint64, slotTxSnapshot *types.TxSlotSnapshot) {
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

// Other txs handling

func (b *txSnapshotsBuilder) getTxPoolSnapshot() *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)
	if txPoolSnapshot == nil {
		return types.NewTxPoolSnapshot()
	}

	return txPoolSnapshot
}

func (b *txSnapshotsBuilder) updateTxPoolSnapshot(txs []*types.Transaction, bz []byte, env *environment) *types.TxPoolSnapshot {
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

	txPoolSnapshot.GasUsed = env.header.GasLimit - env.gasPool.Gas()
	txPoolSnapshot.BytesLength = uint64(len(bz))

	rawdb.WriteTxPoolSnapshot(b.db, txPoolSnapshot)

	return txPoolSnapshot
}

// ProposeFromTxPoolSnapshot proposes txs from the tx pool snapshot.
func (b *txSnapshotsBuilder) proposeFromTxPoolSnapshot() *types.TxPoolSnapshot {
	b.txPoolSnapshotMu.Lock()
	defer b.txPoolSnapshotMu.Unlock()

	txPoolSnapshot := rawdb.ReadTxPoolSnapshot(b.db)

	// Short lived cache to speed up the lookup
	loadProposedInCache(b, txPoolSnapshot)

	// Do not reset the 'NewTxs' field to handle the case of 'driver' failures.
	// Each 'propose' call will accumulate until the 'driver' finally executes
	// and resets the state.
	// txPoolSnapshot.NewTxs = []*types.Transaction{}

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

func (b *txSnapshotsBuilder) resetTxPoolSnapshot() {
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

func loadPendingInCache(b *txSnapshotsBuilder, txPoolSnapshot *types.TxPoolSnapshot) {
	if txPoolSnapshot == nil || (txPoolSnapshot == types.NewTxPoolSnapshot()) || b.pendingTxCache == nil {
		b.pendingTxCache = initTxCache()
		return
	}

	for _, tx := range txPoolSnapshot.PendingTxs {
		b.pendingTxCache[tx.Hash()] = true
	}
}

func loadProposedInCache(b *txSnapshotsBuilder, txPoolSnapshot *types.TxPoolSnapshot) {
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

// resetAllTxSnapshots resets all tx snapshots.
func (b *txSnapshotsBuilder) resetAllTxSnapshots() {
	for i := 0; i < common.EpochLength; i++ {
		b.resetTxSlotSnapshot(uint64(i))
	}
	log.Warn("All tx slot snapshots have been reset")
	b.resetTxPoolSnapshot()
	log.Warn("Tx pool snapshot has been reset")
}
