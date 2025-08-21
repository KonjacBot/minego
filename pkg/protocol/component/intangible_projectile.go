package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type IntangibleProjectile struct {
	Empty nbt.RawMessage `mc:"NBT"` // Always empty compound tag
}

func (*IntangibleProjectile) Type() slot.ComponentID {
	return 19
}

func (*IntangibleProjectile) ID() string {
	return "minecraft:intangible_projectile"
}
