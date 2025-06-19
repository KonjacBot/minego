package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type Enchantments struct {
	Enchantments []Enchantment
}

type Enchantment struct {
	TypeID int32 `mc:"VarInt"`
	Level  int32 `mc:"VarInt"`
}

func (*Enchantments) Type() slot.ComponentID {
	return 10
}

func (*Enchantments) ID() string {
	return "minecraft:enchantments"
}
