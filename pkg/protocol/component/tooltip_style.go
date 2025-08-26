package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type TooltipStyle struct {
	Style pk.Identifier
}

func (*TooltipStyle) Type() slot.ComponentID {
	return 31
}

func (*TooltipStyle) ID() string {
	return "minecraft:tooltip_style"
}
