package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type WolfSoundVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*WolfSoundVariant) Type() slot.ComponentID {
	return 74
}

func (*WolfSoundVariant) ID() string {
	return "minecraft:wolf/sound_variant"
}
