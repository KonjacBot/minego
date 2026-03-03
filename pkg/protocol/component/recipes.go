package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
