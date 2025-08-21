package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Food struct {
	Nutrition          int32 `mc:"VarInt"`
	SaturationModifier float32
	CanAlwaysEat       bool
}

func (*Food) Type() slot.ComponentID {
	return 20
}

func (*Food) ID() string {
	return "minecraft:food"
}
