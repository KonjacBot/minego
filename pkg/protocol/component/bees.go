package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type Bees struct {
	Bees []Bee
}

// Bee 是蜂巢、蜂箱裡的一隻蜜蜂。EntityType、EntityData 是 26.2 的 TypedEntityData:
// 先是實體種類的登錄表編號,再是不含 id 的 NBT。
//
//codec:gen
type Bee struct {
	EntityType     int32          `mc:"VarInt"`
	EntityData     nbt.RawMessage `mc:"NBT"`
	TicksInHive    int32          `mc:"VarInt"`
	MinTicksInHive int32          `mc:"VarInt"`
}

func (*Bees) ID() string {
	return "minecraft:bees"
}
