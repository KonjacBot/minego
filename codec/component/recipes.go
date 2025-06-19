package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type Recipes struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*Recipes) Type() slot.ComponentID {
	return 57
}

func (*Recipes) ID() string {
	return "minecraft:recipes"
}
