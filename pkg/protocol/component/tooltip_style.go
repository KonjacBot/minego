package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type TooltipStyle struct {
	Style pk.Identifier
}

func (*TooltipStyle) ID() string {
	return "minecraft:tooltip_style"
}
