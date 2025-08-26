package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
