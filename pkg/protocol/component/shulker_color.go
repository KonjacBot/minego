package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type ShulkerColor struct {
	Color DyeColor
}

func (*ShulkerColor) Type() slot.ComponentID {
	return 95
}

func (*ShulkerColor) ID() string {
	return "minecraft:shulker/color"
}
