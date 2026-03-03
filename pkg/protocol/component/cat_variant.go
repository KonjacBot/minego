package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
