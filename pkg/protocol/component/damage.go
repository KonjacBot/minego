package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type Damage struct {
	Damage int32 `mc:"VarInt"`
}

func (*Damage) Type() slot.ComponentID {
	return 3
}

func (*Damage) ID() string {
	return "minecraft:damage"
}
