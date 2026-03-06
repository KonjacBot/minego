package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type BreakSound struct {
	SoundData packet.OptID[SoundEvent, *SoundEvent]
}

func (*BreakSound) ID() string {
	return "minecraft:break_sound"
}
