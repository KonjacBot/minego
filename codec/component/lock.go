package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/nbt"
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
