package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
