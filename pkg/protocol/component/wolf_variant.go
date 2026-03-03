package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type WolfVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*WolfVariant) Type() slot.ComponentID {
	return 73
}

func (*WolfVariant) ID() string {
	return "minecraft:wolf/variant"
}
