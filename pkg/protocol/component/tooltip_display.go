package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TooltipDisplay struct {
	HideTooltip      bool
	HiddenComponents []int32 `mc:"VarInt"`
}

func (*TooltipDisplay) Type() slot.ComponentID {
	return 15
}

func (*TooltipDisplay) ID() string {
	return "minecraft:tooltip_display"
}
