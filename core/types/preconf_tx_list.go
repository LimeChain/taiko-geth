package types

// Respresents a snapshot of the tx pool state, used for
// building pre-built tx lists.

type TxPoolSnapshot struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions

	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions

	// TODO(limechain): take into account the anchor tx gas usage and bytes
	EstimatedGasUsed uint64
	BytesLength      uint64
}
