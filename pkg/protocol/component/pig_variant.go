package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type PigVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*PigVariant) Type() slot.ComponentID {
	return 84
}

func (*PigVariant) ID() string {
	return "minecraft:pig/variant"
}
