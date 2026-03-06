package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type Recipes struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*Recipes) ID() string {
	return "minecraft:recipes"
}
