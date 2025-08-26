package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
