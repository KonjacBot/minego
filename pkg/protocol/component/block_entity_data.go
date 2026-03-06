package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type BlockEntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*BlockEntityData) ID() string {
	return "minecraft:block_entity_data"
}
