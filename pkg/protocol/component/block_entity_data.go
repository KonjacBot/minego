package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/nbt"
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
