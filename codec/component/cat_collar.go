package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type CatCollar struct {
	Color DyeColor
}

func (*CatCollar) Type() slot.ComponentID {
	return 93
}

func (*CatCollar) ID() string {
	return "minecraft:cat/collar"
}
