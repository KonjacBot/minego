package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type MapDecorations struct {
	Data nbt.RawMessage `mc:"NBT"` // Always a Compound Tag
}

func (*MapDecorations) Type() slot.ComponentID {
	return 38
}

func (*MapDecorations) ID() string {
	return "minecraft:map_decorations"
}
