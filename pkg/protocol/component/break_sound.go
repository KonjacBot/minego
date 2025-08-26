package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
