package component

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type ParrotVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*ParrotVariant) Type() slot.ComponentID {
	return 78
}

func (*ParrotVariant) ID() string {
	return "minecraft:parrot/variant"
}
