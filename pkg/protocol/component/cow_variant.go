package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type CowVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*CowVariant) Type() slot.ComponentID {
	return 85
}

func (*CowVariant) ID() string {
	return "minecraft:cow/variant"
}
