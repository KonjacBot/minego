package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
