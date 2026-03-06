package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type EntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*EntityData) ID() string {
	return "minecraft:entity_data"
}
