package common

var (
	// TODO(limechain): pass the genesis timestamp as flag or some config
	genesisTimestamp = uint64(1724415819) // genesis timestamp for the L1 chain
)

var (
	epochLength  = uint64(32) // number of slots per epoch
	slotDuration = uint64(12) // duration of each slot in seconds
)

func CurrentSlotAndEpoch(now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	currentSlot := uint64(elapsedTime) / slotDuration
	currentEpoch := currentSlot / epochLength
	return currentSlot, currentEpoch
}
