package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type CustomData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*CustomData) Type() slot.ComponentID {
	return 0
}

func (*CustomData) ID() string {
	return "minecraft:custom_data"
}
