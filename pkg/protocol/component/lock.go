package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
