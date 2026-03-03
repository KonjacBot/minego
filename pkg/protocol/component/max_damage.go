package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type MaxDamage struct {
	Damage int32 `mc:"VarInt"`
}

func (*MaxDamage) Type() slot.ComponentID {
	return 2
}

func (*MaxDamage) ID() string {
	return "minecraft:max_damage"
}
