package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/nbt"
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

func (*Bees) Type() slot.ComponentID {
	return 68
}

func (*Bees) ID() string {
	return "minecraft:bees"
}
