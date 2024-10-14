package common

// CHANGE(limechain):

var (
	EpochLength  = 32 // number of slots per epoch
	SlotDuration = 12 // duration of each slot in seconds
)

const (
	SlotAcceptanceTimeframe = uint64(2) // seconds
)

func HeadSlotStartEndTime(genesisTimestamp uint64, now int64) (uint64, uint64) {
	currentSlot, _ := HeadSlotAndEpoch(genesisTimestamp, now)
	slotStartTime := genesisTimestamp + currentSlot*uint64(SlotDuration)
	slotEndTime := slotStartTime + uint64(SlotDuration)
	return slotStartTime, slotEndTime
}

func HeadSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	currentSlot := uint64(elapsedTime) / uint64(SlotDuration)
	currentEpoch := currentSlot / uint64(EpochLength)
	return currentSlot, currentEpoch
}

func SlotIndex(slot uint64) uint64 {
	return slot % uint64(EpochLength)
}
