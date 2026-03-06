package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type BucketEntityData struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*BucketEntityData) ID() string {
	return "minecraft:bucket_entity_data"
}
