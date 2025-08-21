package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type CatVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*CatVariant) Type() slot.ComponentID {
	return 92
}

func (*CatVariant) ID() string {
	return "minecraft:cat/variant"
}
