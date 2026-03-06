package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type CustomData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*CustomData) ID() string {
	return "minecraft:custom_data"
}
