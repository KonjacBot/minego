package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
