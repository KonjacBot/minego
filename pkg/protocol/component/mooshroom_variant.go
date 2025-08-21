package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type MooshroomVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*MooshroomVariant) Type() slot.ComponentID {
	return 82
}

func (*MooshroomVariant) ID() string {
	return "minecraft:mooshroom/variant"
}
