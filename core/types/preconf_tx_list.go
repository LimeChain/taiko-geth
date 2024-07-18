package types

// PreBuiltTxList is a pre-built transaction list based on the latest chain state,
// with estimated gas used / bytes.
type PreBuiltTxList struct {
	TxList           Transactions
	EstimatedGasUsed uint64
	BytesLength      uint64
}

type TxListState struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions

	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions

	// TODO: take into account the anchor tx gas usage
	EstimatedGasUsed uint64
	BytesLength      uint64
}
