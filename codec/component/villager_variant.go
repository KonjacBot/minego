package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type VillagerVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*VillagerVariant) Type() slot.ComponentID {
	return 72
}

func (*VillagerVariant) ID() string {
	return "minecraft:villager/variant"
}
