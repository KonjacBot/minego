package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type ShulkerColor struct {
	Color DyeColor
}

func (*ShulkerColor) Type() slot.ComponentID {
	return 95
}

func (*ShulkerColor) ID() string {
	return "minecraft:shulker/color"
}
