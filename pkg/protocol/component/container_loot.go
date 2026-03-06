package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type ContainerLoot struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*ContainerLoot) ID() string {
	return "minecraft:container_loot"
}
