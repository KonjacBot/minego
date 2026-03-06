package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type Bees struct {
	Bees []Bee
}

//codec:gen
type Bee struct {
	EntityData     nbt.RawMessage `mc:"NBT"`
	TicksInHive    int32          `mc:"VarInt"`
	MinTicksInHive int32          `mc:"VarInt"`
}

func (*Bees) ID() string {
	return "minecraft:bees"
}
