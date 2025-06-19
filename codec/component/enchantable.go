package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type Enchantable struct {
	Value int32 `mc:"VarInt"`
}

func (*Enchantable) Type() slot.ComponentID {
	return 27
}

func (*Enchantable) ID() string {
	return "minecraft:enchantable"
}
