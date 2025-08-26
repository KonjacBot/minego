package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type BlockEntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*BlockEntityData) Type() slot.ComponentID {
	return 51
}

func (*BlockEntityData) ID() string {
	return "minecraft:block_entity_data"
}
