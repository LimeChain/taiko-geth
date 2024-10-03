package common

// CHANGE(limechain):

const (
	SlotsOffsetInAdvance = 2 // value imposed by the other infra components
)

var (
	EpochLength  = 32 // number of slots per epoch
	SlotDuration = 12 // duration of each slot in seconds
)

func CurrentSlotStartEndTime(genesisTimestamp uint64, now int64) (uint64, uint64) {
	currentSlot, _ := CurrentSlotAndEpoch(genesisTimestamp, now)
	slotStartTime := genesisTimestamp + currentSlot*uint64(SlotDuration)
	slotEndTime := slotStartTime + uint64(SlotDuration)
	return slotStartTime, slotEndTime
}

func CurrentSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	currentSlot := uint64(elapsedTime) / uint64(SlotDuration)
	currentEpoch := currentSlot / uint64(EpochLength)
	return currentSlot, currentEpoch
}

func SlotIndex(slot uint64) uint64 {
	return slot % uint64(EpochLength)
}
