package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type PotDecorations struct {
	Decorations []int32 `mc:"PrefixedArray"`
}

func (*PotDecorations) Type() slot.ComponentID {
	return 65
}

func (*PotDecorations) ID() string {
	return "minecraft:pot_decorations"
}
