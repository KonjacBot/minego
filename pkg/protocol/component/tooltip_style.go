package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type TooltipStyle struct {
	Style wire.Identifier
}

func (*TooltipStyle) ID() string {
	return "minecraft:tooltip_style"
}
