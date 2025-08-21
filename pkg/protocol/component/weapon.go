package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Weapon struct {
	DamagePerAttack    int32   `mc:"VarInt"`
	DisableBlockingFor float32 // In seconds
}

func (*Weapon) Type() slot.ComponentID {
	return 26
}

func (*Weapon) ID() string {
	return "minecraft:weapon"
}
