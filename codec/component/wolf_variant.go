package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type WolfVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*WolfVariant) Type() slot.ComponentID {
	return 73
}

func (*WolfVariant) ID() string {
	return "minecraft:wolf/variant"
}
