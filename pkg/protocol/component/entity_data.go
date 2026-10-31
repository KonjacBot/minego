package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

// EntityData 是生怪蛋、盔甲座這類物品帶的實體資料(26.2 的 TypedEntityData):
// 先是實體種類的登錄表編號,再是不含 id 的 NBT。
//
//codec:gen
type EntityData struct {
	Type int32          `mc:"VarInt"`
	Data nbt.RawMessage `mc:"NBT"`
}

func (*EntityData) ID() string {
	return "minecraft:entity_data"
}
