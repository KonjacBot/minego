package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type FrogVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*FrogVariant) Type() slot.ComponentID {
	return 87
}

func (*FrogVariant) ID() string {
	return "minecraft:frog/variant"
}
