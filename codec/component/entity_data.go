package component

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type EntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*EntityData) Type() slot.ComponentID {
	return 49
}

func (*EntityData) ID() string {
	return "minecraft:entity_data"
}
