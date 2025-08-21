package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TropicalFishBaseColor struct {
	Color DyeColor
}

func (*TropicalFishBaseColor) Type() slot.ComponentID {
	return 80
}

func (*TropicalFishBaseColor) ID() string {
	return "minecraft:tropical_fish/base_color"
}
