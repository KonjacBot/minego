package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

// BlockEntityData 是方塊物品帶的方塊實體資料(26.2 的 TypedEntityData):
// 先是方塊實體種類的登錄表編號,再是不含 id 的 NBT。
//
//codec:gen
type BlockEntityData struct {
	Type int32          `mc:"VarInt"`
	Data nbt.RawMessage `mc:"NBT"`
}

func (*BlockEntityData) ID() string {
	return "minecraft:block_entity_data"
}
