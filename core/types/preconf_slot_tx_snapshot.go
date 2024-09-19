package types

// CHANGE(limechain):

// SlotTxSnapshot represents a snapshot of the txs being processed.
type SlotTxSnapshot struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions
	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions
	// Keeps track of the block space constraints per slot, in the current epoch prior actual execution
	GasUsed     uint64
	BytesLength uint64
}
