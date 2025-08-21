package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type AxolotlVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*AxolotlVariant) Type() slot.ComponentID {
	return 91
}

func (*AxolotlVariant) ID() string {
	return "minecraft:axolotl/variant"
}
