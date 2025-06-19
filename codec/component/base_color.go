package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type BaseColor struct {
	Color DyeColor
}

func (*BaseColor) Type() slot.ComponentID {
	return 64
}

func (*BaseColor) ID() string {
	return "minecraft:base_color"
}
