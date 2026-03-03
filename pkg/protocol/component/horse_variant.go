package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type HorseVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*HorseVariant) Type() slot.ComponentID {
	return 88
}

func (*HorseVariant) ID() string {
	return "minecraft:horse/variant"
}
