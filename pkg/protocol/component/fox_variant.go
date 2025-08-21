package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type FoxVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*FoxVariant) Type() slot.ComponentID {
	return 76
}

func (*FoxVariant) ID() string {
	return "minecraft:fox/variant"
}
