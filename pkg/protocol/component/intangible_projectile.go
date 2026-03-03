package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
