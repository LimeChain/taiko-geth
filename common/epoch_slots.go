package common

var (
	epochLength  = uint64(32) // number of slots per epoch
	slotDuration = uint64(12) // duration of each slot in seconds
)

func CurrentSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	currentSlot := uint64(elapsedTime) / slotDuration
	currentEpoch := currentSlot / epochLength
	return currentSlot, currentEpoch
}
