package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

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
