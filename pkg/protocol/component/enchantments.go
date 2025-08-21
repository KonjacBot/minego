package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Enchantments struct {
	Enchantments []Enchantment
}

//codec:gen
type Enchantment struct {
	Type  int32 `mc:"VarInt"`
	Level int32 `mc:"VarInt"`
}

func (*Enchantments) Type() slot.ComponentID {
	return 10
}

func (*Enchantments) ID() string {
	return "minecraft:enchantments"
}
