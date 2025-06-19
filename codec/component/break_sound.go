package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
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
