package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/nbt"
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
