package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TropicalFishPattern struct {
	Pattern int32 `mc:"VarInt"`
}

func (*TropicalFishPattern) Type() slot.ComponentID {
	return 79
}

func (*TropicalFishPattern) ID() string {
	return "minecraft:tropical_fish/pattern"
}
