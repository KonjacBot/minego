package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type CreativeSlotLock struct {
	// no fields
}

func (*CreativeSlotLock) Type() slot.ComponentID {
	return 17
}

func (*CreativeSlotLock) ID() string {
	return "minecraft:creative_slot_lock"
}
