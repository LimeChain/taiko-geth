package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

// CHANGE(limechain): keep track of the block constraints per slot.

type BlockConstraints struct {
	EstimatedGasUsed uint64
	BytesLength      uint64
}

type PerSlotConstraints struct {
	// Keeps track of the block space constraints per slot,
	// in the current epoch prior actual execution
	SlotConstraints [32]*BlockConstraints

	// Keeps track of the current usage during execution (might be removed at some point)
	Total BlockConstraints
}

func NewPerSlotConstraints() *PerSlotConstraints {
	return &PerSlotConstraints{
		Total: BlockConstraints{},
		SlotConstraints: [32]*BlockConstraints{
			{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		},
	}
}

// Get returns constraints for a given slot.
func (s *PerSlotConstraints) Get(slot uint64) (*BlockConstraints, error) {
	index := slot % 32
	constraints := s.SlotConstraints[index]
	if constraints == nil {
		return &BlockConstraints{}, fmt.Errorf("Block constraints for slot %d not found", slot)
	}
	return constraints, nil
}

// Update updates the constraints for a given slot.
func (s *PerSlotConstraints) Update(slot uint64, gasUsed uint64, bytesLength uint64) error {
	index := slot % 32
	constraints := s.SlotConstraints[index]
	if constraints == nil {
		return fmt.Errorf("Block constraints for slot %d not found", slot)
	}
	constraints.EstimatedGasUsed = gasUsed
	constraints.BytesLength = bytesLength
	return nil
}

// Reset resets past constraints prior to the current slot.
func (s *PerSlotConstraints) Reset(currentSlot uint64) {
	for i := uint64(0); i < currentSlot; i++ {
		index := i % 32
		s.SlotConstraints[index] = &BlockConstraints{}
	}
	log.Info("Past slot constraints were reset")
}
