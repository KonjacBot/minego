package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type Lock struct {
	Key nbt.RawMessage `mc:"NBT"`
}

func (*Lock) ID() string {
	return "minecraft:lock"
}
