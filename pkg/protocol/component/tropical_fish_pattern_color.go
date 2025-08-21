package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TropicalFishPatternColor struct {
	Color DyeColor
}

func (*TropicalFishPatternColor) Type() slot.ComponentID {
	return 81
}

func (*TropicalFishPatternColor) ID() string {
	return "minecraft:tropical_fish/pattern_color"
}
