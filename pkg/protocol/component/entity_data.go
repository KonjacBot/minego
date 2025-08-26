package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
