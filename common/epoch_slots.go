package common

// CHANGE(limechain):

const (
	SlotsOffsetInAdvance = 3 // value requirement, imposed by the other infra components
)

var (
	EpochLength  = 32 // number of slots per epoch
	SlotDuration = 12 // duration of each slot in seconds
)

func CurrentSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	currentSlot := uint64(elapsedTime) / uint64(SlotDuration)
	currentEpoch := currentSlot / uint64(EpochLength)
	return currentSlot, currentEpoch
}

func SlotIndex(slot uint64) uint64 {
	return slot % uint64(EpochLength)
}
