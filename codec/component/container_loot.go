package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type ContainerLoot struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*ContainerLoot) Type() slot.ComponentID {
	return 70
}

func (*ContainerLoot) ID() string {
	return "minecraft:container_loot"
}
