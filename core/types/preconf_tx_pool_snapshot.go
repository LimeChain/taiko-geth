package types

import (
	"fmt"
)

// CHANGE(limechain): keep track of the current txs and block constraints per slot.

type BlockConstraints struct {
	EstimatedGasUsed uint64
	BytesLength      uint64
}

// TxPoolSnapshot represents a snapshot of the txs being processed.
type TxPoolSnapshot struct {
	// Keeps track of current pending txs
	PendingTxs Transactions
	// Keeps track of current pending txs that heve been proposed
	ProposedTxs Transactions
	// New txs that are ready to be proposed, kept until actually executed
	NewTxs Transactions

	// Keeps track of the current usage during execution
	BlockConstraints

	// Keeps track of the block space constraints per slot,
	// in the current epoch prior actual execution
	PerSlotConstraints [32]*BlockConstraints
}

func NewTxPoolSnapshot() *TxPoolSnapshot {
	return &TxPoolSnapshot{
		PendingTxs:  Transactions{},
		ProposedTxs: Transactions{},
		NewTxs:      Transactions{},

		BlockConstraints: BlockConstraints{},

		PerSlotConstraints: [32]*BlockConstraints{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		},
	}
}

// GetConstraints returns constraints for a given slot.
func (s *TxPoolSnapshot) GetConstraints(slot uint64) (*BlockConstraints, error) {
	index := slot % 32
	constraints := s.PerSlotConstraints[index]
	if constraints == nil {
		return &BlockConstraints{}, fmt.Errorf("Block constraints for slot %d not found", slot)
	}
	return constraints, nil
}

// UpdateConstraints updates the constraints for a given slot.
func (s *TxPoolSnapshot) UpdateConstraints(slot uint64, gasUsed uint64, bytesLength uint64) error {
	index := slot % 32
	constraints := s.PerSlotConstraints[index]
	if constraints == nil {
		return fmt.Errorf("Block constraints for slot %d not found", slot)
	}
	constraints.EstimatedGasUsed = gasUsed
	constraints.BytesLength = bytesLength
	return nil
}

// ResetPastConstraints resets past constraints based on the current slot.
func (s *TxPoolSnapshot) ResetPastConstraints(currentSlot uint64) error {
	index := currentSlot % 32
	for i := index - 1; i > 0; i-- {
		s.PerSlotConstraints[i] = &BlockConstraints{}
	}
	return nil
}
