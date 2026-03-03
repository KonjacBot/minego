package component

import (
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type BreakSound struct {
	SoundData packet.OptID[SoundEvent, *SoundEvent]
}

func (*BreakSound) Type() slot.ComponentID {
	return 71
}

func (*BreakSound) ID() string {
	return "minecraft:break_sound"
}
