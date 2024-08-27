package types

// CHANGE(limechain):

// TxPoolSnapshot represents a snapshot of the txs being processed.
type TxPoolSnapshot struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions
	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions
}
