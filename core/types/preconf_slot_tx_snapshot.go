package types

// CHANGE(limechain):

// Snapshot of non-preconf txs from the tx pool, currently being processed.
type TxPoolSnapshot struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions
	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions
	// Keeps track of the block space constraints.
	GasUsed     uint64
	BytesLength uint64
}

func NewTxPoolSnapshot() *TxPoolSnapshot {
	return &TxPoolSnapshot{
		PendingTxs:  make(Transactions, 0),
		ProposedTxs: make(Transactions, 0),
		NewTxs:      make(Transactions, 0),
	}
}

// Snapshot of preconf txs prepeared in advance for each slot,
// there is no case where proposer runs during preparation of the slot snapshot.
type TxSlotSnapshot struct {
	SlotIndex   uint64
	Txs         Transactions
	GasUsed     uint64
	BytesLength uint64
}

func NewTxSlotSnapshot() *TxSlotSnapshot {
	return &TxSlotSnapshot{
		Txs: make(Transactions, 0),
	}
}
