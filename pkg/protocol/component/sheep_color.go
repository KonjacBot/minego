package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type SheepColor struct {
	Color DyeColor
}

func (*SheepColor) Type() slot.ComponentID {
	return 94
}

func (*SheepColor) ID() string {
	return "minecraft:sheep/color"
}
