package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type IntangibleProjectile struct {
	Empty nbt.RawMessage `mc:"NBT"` // Always empty compound tag
}

func (*IntangibleProjectile) ID() string {
	return "minecraft:intangible_projectile"
}
