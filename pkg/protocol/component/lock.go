package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type Lock struct {
	Key nbt.RawMessage `mc:"NBT"`
}

func (*Lock) Type() slot.ComponentID {
	return 69
}

func (*Lock) ID() string {
	return "minecraft:lock"
}
