package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type MapDecorations struct {
	Data nbt.RawMessage `mc:"NBT"` // Always a Compound Tag
}

func (*MapDecorations) ID() string {
	return "minecraft:map_decorations"
}
