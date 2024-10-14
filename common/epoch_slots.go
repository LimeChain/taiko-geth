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
	headSlot, _ := HeadSlotAndEpoch(genesisTimestamp, now)
	slotStartTime := genesisTimestamp + headSlot*uint64(SlotDuration)
	slotEndTime := slotStartTime + uint64(SlotDuration)
	return slotStartTime, slotEndTime
}

func HeadSlotAndEpoch(genesisTimestamp uint64, now int64) (uint64, uint64) {
	elapsedTime := uint64(now) - genesisTimestamp
	headSlot := uint64(elapsedTime) / uint64(SlotDuration)
	headEpoch := headSlot / uint64(EpochLength)
	return headSlot, headEpoch
}

func SlotIndex(slot uint64) uint64 {
	return slot % uint64(EpochLength)
}
