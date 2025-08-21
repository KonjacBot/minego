package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	pk "github.com/Tnze/go-mc/net/packet"
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
